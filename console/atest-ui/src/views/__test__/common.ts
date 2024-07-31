/*
Copyright 2024 API Testing Authors.

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

export function SetupStorage() {
    const localStorageMock = (() => {
        let store = new Map<string, string>;

        return {
            getItem(key: string) {
                return store.get(key) || null;
            },
            setItem(key: string, value: string) {
                store.set(key, value);
            },
            removeItem(key: string) {
                store.delete(key)
            },
            clear() {
                store.clear();
            }
        };
    })();

    Object.defineProperty(global, 'sessionStorage', {
        value: localStorageMock
    });
    Object.defineProperty(global, 'localStorage', {
        value: localStorageMock
    });
    Object.defineProperty(global, 'navigator', {
        value: {
            language: 'en'
        }
    });
}
