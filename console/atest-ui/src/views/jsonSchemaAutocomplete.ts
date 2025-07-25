import type { Completion } from "@codemirror/autocomplete";
import { CompletionContext } from "@codemirror/autocomplete";
import type { CompletionResult } from "@codemirror/autocomplete";
import { syntaxTree } from "@codemirror/language";
import type { JSONSchema7 } from "json-schema";
import Ajv from "ajv";
import { parser } from "@lezer/yaml";
import * as yaml from "yaml";
import { yamlLanguage } from "@codemirror/lang-yaml";

/**
 * 基于 JSON Schema 的 yaml 自动补全函数
 * @param schema JSONSchema7 对象
 * @returns CodeMirror 自动补全源
 */
export function yamlSchemaAutocomplete(schema: JSONSchema7) {
    return (context: CompletionContext): CompletionResult | null => {
        const tree = syntaxTree(context.state);
        const node = tree.resolveInner(context.pos, -1);

        // 解析 YAML 路径
        let path: (string | number)[] = [];
        let currentNode = node;
        while (currentNode) {
            if (currentNode.type.name === "Property") {
                const keyNode = currentNode.getChild("Key");
                if (keyNode) {
                    const key = context.state.sliceDoc(keyNode.from, keyNode.to);
                    path.unshift(key.replace(/^['"]|['"]$/g, ''));
                }
            } else if (currentNode.type.name === "Item") {
                let parent = currentNode.parent;
                while (parent && parent.type.name !== "Property") {
                    parent = parent.parent;
                }
                if (parent) {
                    const keyNode = parent.getChild("Key");
                    if (keyNode) {
                        const key = context.state.sliceDoc(keyNode.from, keyNode.to);
                        path.unshift(key.replace(/^['"]|['"]$/g, ''));
                    }
                }
                path.push("[]");
            }
            currentNode = currentNode.parent;
        }

        // 根据路径定位 schema
        let currentSchema: any = schema;
        for (const key of path) {
            if (key === "[]") {
                if (currentSchema.items) {
                    currentSchema = currentSchema.items;
                } else {
                    return null;
                }
            } else if (currentSchema.properties && currentSchema.properties[key]) {
                currentSchema = currentSchema.properties[key];
            } else {
                return null;
            }
        }

        // 生成补全项
        const options: Completion[] = [];
        if (currentSchema.properties) {
            for (const [key, prop] of Object.entries(currentSchema.properties)) {
                options.push({
                    label: key,
                    type: "property",
                    detail: prop.description || `type: ${prop.type}`,
                    apply: key + ": "
                });
            }
        } else if (currentSchema.enum) {
            for (const value of currentSchema.enum) {
                options.push({
                    label: String(value),
                    type: "value",
                    detail: currentSchema.description
                });
            }
        } else if (currentSchema.type === "boolean") {
            options.push(
                { label: "true", type: "value" },
                { label: "false", type: "value" }
            );
        } else if (currentSchema.type === "array" && currentSchema.items) {
            options.push({
                label: "- ",
                type: "array",
                detail: "Array item"
            });
        }

        return options.length > 0 ? {
            from: context.pos,
            options,
            validFor: /^[\w$-]*$/
        } : null;
    };
}

export function jsonSchemaCompletions(schema: JSONSchema7) {
    return (context: any) => {
        const tree = syntaxTree(context.state);
        const node = tree.resolveInner(context.pos, -1);

        // 获取当前路径
        let path: (string | number)[] = [];
        let currentNode = node;
        while (currentNode) {
            if (currentNode.type.name === "Property") {
                const keyNode = currentNode.getChild("Key");
                if (keyNode) {
                    const key = context.state.sliceDoc(keyNode.from, keyNode.to);
                    path.unshift(key.replace(/^['"]|['"]$/g, '')); // 移除引号
                }
            } else if (currentNode.type.name === "Item") {
                // YAML 数组元素
                // 获取父节点的 key
                let parent = currentNode.parent;
                while (parent && parent.type.name !== "Property") {
                    parent = parent.parent;
                }
                if (parent) {
                    const keyNode = parent.getChild("Key");
                    if (keyNode) {
                        const key = context.state.sliceDoc(keyNode.from, keyNode.to);
                        path.unshift(key.replace(/^['"]|['"]$/g, ''));
                    }
                }
                // 标记为数组项
                path.push("[]");
            }
            currentNode = currentNode.parent;
        }

        // 根据路径获取当前 Schema 定义
        let currentSchema: any = schema;
        for (const key of path) {
            if (key === "[]") {
                if (currentSchema.items) {
                    currentSchema = currentSchema.items;
                } else {
                    return null;
                }
            } else if (currentSchema.properties && currentSchema.properties[key]) {
                currentSchema = currentSchema.properties[key];
            } else {
                return null; // 路径不匹配
            }
        }

        // 生成补全选项
        const options: Completion[] = [];
        if (currentSchema.properties) {
            // 对象属性补全
            for (const [key, prop] of Object.entries(currentSchema.properties)) {
                options.push({
                    label: key,
                    type: "property",
                    detail: prop.description || `type: ${prop.type}`,
                    apply: key + ": "
                });
            }
        } else if (currentSchema.enum) {
            // 枚举值补全
            for (const value of currentSchema.enum) {
                options.push({
                    label: String(value),
                    type: "value",
                    detail: currentSchema.description
                });
            }
        } else if (currentSchema.type === "boolean") {
            // 布尔值补全
            options.push(
                { label: "true", type: "value" },
                { label: "false", type: "value" }
            );
        } else if (currentSchema.type === "array" && currentSchema.items) {
            // 数组类型补全
            options.push({
                label: "- ",
                type: "array",
                detail: "Array item"
            });
        }

        return {
            from: context.pos,
            options,
            validFor: /^[\w$-]*$/
        };
    };
}

// 根据 JSON Schema 生成补全建议
export function jsonSchemaCompleter(schema: JSONSchema7) {
  return (context: CompletionContext): CompletionResult | null => {
    const node = syntaxTree(context.state).resolveInner(context.pos, -1);
    const text = context.state.sliceDoc(0, context.pos);
    
    // 检查是否在键位置（属性名）
    if (node.name === "PropertyName") {
      return completePropertyKeys(schema, context);
    }
    
    // 检查是否在值位置（字符串/枚举值）
    if (node.parent?.name === "Property") {
      const keyNode = node.parent.firstChild;
      if (keyNode) {
        const key = context.state.sliceDoc(keyNode.from + 1, keyNode.to - 1); // 去掉引号
        return completePropertyValues(schema, key, context);
      }
    }
    
    return null;
  };
}

// 补全属性名
function completePropertyKeys(
  schema: JSONSchema7,
  context: CompletionContext
): CompletionResult | null {
  const resolvedSchema = resolveSchema(schema) as JSONSchema7;
  const properties = resolvedSchema.properties || {};
  const required = resolvedSchema.required || [];

  const options: Completion[] = Object.entries(properties).map(([key, def]) => ({
    label: key,
    type: "property",
    detail: def.description || `${def.type}`,
    boost: required.includes(key) ? 1 : 0, // 高亮必填字段
  }));

  return {
    from: context.pos,
    options,
    validFor: /^[\w$]*$/
  };
}

// 补全属性值（支持枚举和类型）
function completePropertyValues(
  schema: JSONSchema7,
  key: string,
  context: CompletionContext
): CompletionResult | null {
  const resolvedSchema = resolveSchema(schema) as JSONSchema7;
  const property = resolvedSchema.properties?.[key];
  if (!property) return null;

  let options: Completion[] = [];

  // 枚举值补全
  if (property.enum) {
    options = property.enum.map(value => ({
      label: JSON.stringify(value),
      type: "value",
    }));
  }
  // 布尔值补全
  else if (property.type === "boolean") {
    options = ["true", "false"].map(value => ({
      label: value,
      type: "value",
    }));
  }
  // 数组类型提示
  else if (property.type === "array" && property.items) {
    options = [{
      label: "[]",
      type: "array",
      detail: "Array",
    }];
  }

  return options.length > 0 ? {
    from: context.pos,
    options,
  } : null;
}
