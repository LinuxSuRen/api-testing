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

  return pairCompletions(context, templateFuncs.value)
}

const headerFuncs = ref([] as Completion[])
function headerCompletions(context: CompletionContext) {
  if (headerFuncs.value.length == 0) {
    API.PopularHeaders((e) => {
      if (e.data) {
        e.data.forEach((item: any) => {
          headerFuncs.value.push({
            label: item.key,
            type: "text",
          } as Completion)
        })
      }
    })
  }

  return pairCompletions(context, headerFuncs.value)
}

function pairCompletions(context: CompletionContext, source: Completion[]) {
  let word = context.matchBefore(/\w*/) || {
    from: "",
    to: ""
  }
  if (word.from == word.to && !context.explicit)
    return null
  return {
    from: word.from,
    options: source
  }
}

export const NewTemplateLangComplete = (lang: Language) => {
    return new LanguageSupport(lang, lang.data.of(
        {autocomplete: tplFuncCompletions}
    ))
}

export const NewHeaderLangComplete = (lang: Language) => {
    return new LanguageSupport(lang, lang.data.of(
        {autocomplete: headerCompletions}
    ))
}
