/*
Copyright 2024-2025 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import { watch } from 'vue'
import { useMagicKeys } from '@vueuse/core'
import { API } from './net'

function Keys(func: (() => void) | ((k: string) => void), keys: string[]) {
    const magicKeys = useMagicKeys()
    keys.forEach(k => {
        watch(magicKeys[k], (v) => {
            if (v) func(k)
        })
    })
}

interface MagicKey {
    Key?: string
    Keys: string[]
    Func: (() => void) | ((k: string) => void)
    Description?: string
}

const MagicKeyEventName = 'show-key-bindings-dialog'
const AdvancedKeys = (keys: MagicKey[]) => {
    keys.push({
        Keys: ['Ctrl+/'],
        Func: () => {
            const event = new CustomEvent(MagicKeyEventName, {
                detail: keys.map((k) => ({
                    keys: k.Keys.join(', '),
                    description: k.Description ?? 'No description',
                })),
            })
            window.dispatchEvent(event)
        },
        Description: 'Show key bindings dialog',
    })

    keys.forEach((key: MagicKey) => {
        Keys(key.Func, key.Keys)
    })
}

interface KeyBindings {
    name: string
    pages: Page[]
}

interface Page {
    name: string
    bindings: KeyBinding[]
}

interface KeyBinding {
    keys: string[]
    description?: string
    action: string
}

const LoadMagicKeys = (pageName: String, mapping: Map<String, Function>) => {
    console.log(`Loading magic keys for page: ${pageName}`);
    API.GetBinding("default", (data) => {
        const bindings = JSON.parse(data.message) as KeyBindings;
        bindings.pages.forEach((page: Page) => {
            if (page.name === pageName) {
                const keys = [] as MagicKey[]
                page.bindings.forEach((binding: KeyBinding) => {
                    keys.push({
                        Keys: binding.keys,
                        Func: () => {
                            mapping.has(binding.action) ? mapping.get(binding.action)!() : console.warn(`No action found for ${binding.action}`);
                        },
                        Description: binding.description,
                    });
                });
                AdvancedKeys(keys);
                return;
            }
        })
    })
}

export const Magic = {
    Keys, AdvancedKeys, MagicKeyEventName, LoadMagicKeys
}
