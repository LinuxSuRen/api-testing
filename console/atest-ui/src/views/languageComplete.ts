import { Language, LanguageSupport } from "@codemirror/language"
import { CompletionContext } from "@codemirror/autocomplete"
import type { Completion } from "@codemirror/autocomplete"
import { API } from './net'
import { ref } from 'vue'

const templateFuncs = ref([] as Completion[])
function tplFuncCompletions(context: CompletionContext) {
  if (templateFuncs.value.length == 0) {
    API.FunctionsQuery("", "template", (e) => {
      if (e.data) {
        e.data.forEach((item: any) => {
          templateFuncs.value.push({
            label: item.key,
            type: "text",
          } as Completion)
        })
      }
    })
  }

  let word = context.matchBefore(/\w*/) || {
    from: "",
    to: ""
  }
  if (word.from == word.to && !context.explicit)
    return null
  return {
    from: word.from,
    options: templateFuncs.value
  }
}

export const NewLanguageComplete = (lang: Language) => {
    return new LanguageSupport(lang, lang.data.of(
        {autocomplete: tplFuncCompletions}
    ))
}
