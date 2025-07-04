/*
Copyright 2023-2025 API Testing Authors.

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
package render

import (
    "bytes"
    "context"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "io"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
    tests := []struct {
        name   string
        text   string
        ctx    interface{}
        expect string
        verify func(*testing.T, string)
    }{{
        name:   "default",
        text:   `{{default "hello" .Bar}}`,
        ctx:    nil,
        expect: "hello",
    }, {
        name:   "trim",
        text:   `{{trim "   hello    "}}`,
        ctx:    "",
        expect: "hello",
    }, {
        name: "randomKubernetesName",
        text: `{{randomKubernetesName}}`,
        verify: func(t *testing.T, s string) {
            assert.Equal(t, 8, len(s))
        },
    }, {
        name:   "md5",
        text:   `{{md5 "linuxsuren"}}`,
        expect: "b559b80ae1ba1c292d9b3265f265e76a",
    }, {
        name:   "base64",
        text:   `{{base64 "linuxsuren"}}`,
        expect: "bGludXhzdXJlbg==",
    }, {
        name:   "base64Decode",
        text:   `{{base64Decode "bGludXhzdXJlbg=="}}`,
        expect: "linuxsuren",
    }, {
        name:   "base64Decode with error",
        text:   `{{base64Decode "error"}}`,
        expect: "illegal base64 data at input byte 4",
    }, {
        name: "complex",
        text: `{{(index .items 0).name}}?a=a&key={{randomKubernetesName}}`,
        ctx: map[string]interface{}{
            "items": []interface{}{map[string]string{
                "name": "one",
            }, map[string]string{
                "name": "two",
            }},
        },
        verify: func(t *testing.T, s string) {
            assert.Equal(t, 20, len(s), s)
        },
    }, {
        name: "randFloat",
        text: `{{randFloat 1 2}}`,
        verify: func(t *testing.T, s string) {
            assert.NotEmpty(t, s)
        },
    }, {
        name: "randEnum",
        text: `{{randEnum "a" "b" "c"}}`,
        verify: func(t *testing.T, s string) {
            assert.Contains(t, []string{"a", "b", "c"}, s)
        },
    }, {
        name: "randEnumByStr",
        text: `{{randEnumByStr "a,b,c"}}`,
        verify: func(t *testing.T, s string) {
            assert.Contains(t, []string{"a", "b", "c"}, s)
        },
    }, {
        name: "randEnumByJSON",
        text: `{{(randEnumByJSON "[{\"key\":\"a\"},{\"key\":\"b\"},{\"key\":\"c\"}]").key}}`,
        verify: func(t *testing.T, s string) {
            assert.Contains(t, []string{"a", "b", "c"}, s)
        },
    }, {
        name: "randWeightEnum",
        text: `{{randWeightEnum (weightObject 1 "a") (weightObject 2 "b") (weightObject 3 "c")}}`,
        verify: func(t *testing.T, s string) {
            assert.Contains(t, []string{"a", "b", "c"}, s)
        },
    }, {
        name: "randEmail",
        text: `{{randEmail}}`,
        verify: func(t *testing.T, s string) {
            assert.Contains(t, s, "@")
            assert.Contains(t, s, ".com")
        },
    }, {
        name: "sha256 bytes",
        text: `{{sha256sumBytes .data}}`,
        ctx: map[string]interface{}{
            "data": []byte("hello"),
        },
        verify: func(t *testing.T, s string) {
            assert.Equal(t, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", s)
        },
    }, {
        name: "uptimeDate",
        text: `{{uptimeDate}}`,
        verify: func(t *testing.T, s string) {
            assert.NotEmpty(t, s)
        },
    }, {
        name: "uptime",
        text: `{{uptime}}`,
        verify: func(t *testing.T, s string) {
            assert.NotEmpty(t, s)
        },
    }, {
        name: "uptimeSeconds",
        text: `{{uptimeSeconds}}`,
        verify: func(t *testing.T, s string) {
            assert.NotEmpty(t, s)
        },
    }, {
        name: "url encode",
        text: `{{urlEncode "hello world"}}`,
        verify: func(t *testing.T, s string) {
            assert.Equal(t, "hello+world", s)
        },
    }, {
        name: "url decode",
        text: `{{urlDecode "hello+world"}}`,
        verify: func(t *testing.T, s string) {
            assert.Equal(t, "hello world", s)
        },
    }}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Render(tt.name, tt.text, tt.ctx)
            assert.Nil(t, err, err)
            if tt.expect != "" {
                assert.Equal(t, tt.expect, result)
            }
            if tt.verify != nil {
                tt.verify(t, result)
            }
        })
    }
}

func TestRenderThenPrint(t *testing.T) {
    tests := []struct {
        name    string
        tplText string
        ctx     interface{}
        buf     *bytes.Buffer
        expect  string
    }{{
        name:    "simple",
        tplText: `{{max 1 2 3}}`,
        ctx:     nil,
        buf:     new(bytes.Buffer),
        expect:  "3",
    }, {
        name:    "with a map as context",
        tplText: `{{.name}}`,
        ctx:     map[string]string{"name": "linuxsuren"},
        buf:     new(bytes.Buffer),
        expect:  "linuxsuren",
    }}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := RenderThenPrint(tt.name, tt.tplText, tt.ctx, tt.buf)
            assert.NoError(t, err)
            assert.Equal(t, tt.expect, tt.buf.String())
        })
    }
}

func TestFuncGenerator(t *testing.T) {
    tests := []struct {
        name     string
        funcName string
        fields   string
        expect   string
    }{{
        name:     "randomKubernetesName",
        funcName: "randomKubernetesName",
        expect:   `{{randomKubernetesName}}`,
    }, {
        name:     "generateJSONString",
        funcName: "generateJSONString",
        fields:   "name, age",
        expect:   `{{generateJSONString "name" "age"}}`,
    }}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            funcs := GetAdvancedFuncs()
            for _, f := range funcs {
                if f.FuncName == tt.funcName {
                    buf := new(bytes.Buffer)
                    ctx := context.Background()
                    ctx = context.WithValue(ctx, ContextBufferKey, buf)
                    err := f.Generator(ctx, tt.fields)
                    assert.NoError(t, err)
                    assert.Equal(t, tt.expect, buf.String())
                }
            }
        })
    }
}

func TestGoDogGenerator(t *testing.T) {
    tests := []struct {
        name       string
        goDogExper string
        fields     string
        expect     string
    }{{
        name:       "randomKubernetesName",
        goDogExper: `^生成随机字符串，长度 (.*)$`,
        fields:     `3`,
        expect:     `{{randAlpha 3}}`,
    }}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            funcs := GetAdvancedFuncs()
            for _, f := range funcs {
                if f.GoDogExper == tt.goDogExper {
                    buf := new(bytes.Buffer)
                    ctx := context.Background()
                    ctx = context.WithValue(ctx, ContextBufferKey, buf)
                    err := f.Generator(ctx, tt.fields)
                    assert.NoError(t, err)
                    assert.Equal(t, tt.expect, buf.String())
                }
            }
        })
    }
}

func TestGenerateJSONString(t *testing.T) {
    result := generateJSONString([]string{"name", "age"})
    assert.Equal(t, `{"age":"random","name":"random"}`, result)
}

func TestSecret(t *testing.T) {
    SetSecretGetter(nil)
    result, err := Render("", `{{secretValue "pass"}}`, nil)
    assert.NoError(t, err)
    assert.Equal(t, "no secret server", result)

    expected := "password"
    SetSecretGetter(&nonSecretGetter{
        value: expected,
    })
    result, err = Render("", `{{secretValue "pass"}}`, nil)
    assert.NoError(t, err)
    assert.Equal(t, expected, result)

    t.Run("render as reader", func(t *testing.T) {
        reader, err := RenderAsReader("render as reader", "hello", nil)
        assert.NoError(t, err)
        assert.NotNil(t, reader)

        data, err := io.ReadAll(reader)
        assert.NoError(t, err)
        assert.Equal(t, "hello", string(data))
    })
}

func TestRasEncryptWithPublicKey(t *testing.T) {
    // Generate a new RSA key pair
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        t.Fatalf("Failed to generate private key: %v", err)
    }
    publicKey := &privateKey.PublicKey

    // Encode the public key to PEM format
    pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
    if err != nil {
        t.Fatalf("Failed to marshal public key: %v", err)
    }
    pubBytes := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: pubASN1,
    })

    // Encrypt a message using the public key
    message := "hello world"
    encryptedMessage, err := rasEncryptWithPublicKey(message, string(pubBytes))
    if err != nil {
        t.Fatalf("Failed to encrypt message: %v", err)
    }

    // Decrypt the message using the private key
    decodedMessage, err := base64.StdEncoding.DecodeString(encryptedMessage)
    if err != nil {
        t.Fatalf("Failed to decode message: %v", err)
    }
    decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedMessage)
    if err != nil {
        t.Fatalf("Failed to decrypt message: %v", err)
    }

    // Verify the decrypted message
    decryptedMessage := string(decryptedBytes)
    if decryptedMessage != message {
        t.Fatalf("Decrypted message does not match original. Got: %s, want: %s", decryptedMessage, message)
    }
}

func TestFuncUsages(t *testing.T) {
    funcs := []string{"randImage", "randZip"}

    for _, f := range funcs {
        usage := FuncUsage(f)
        assert.NotEmpty(t, usage)
    }
}

func TestGetEngineVersion(t *testing.T) {
    ver := GetEngineVersion()
    assert.Empty(t, ver)
}
