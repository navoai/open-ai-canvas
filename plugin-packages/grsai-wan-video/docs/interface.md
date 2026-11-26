# Grsai Wan Video 接口字段

## 协议身份

- 插件 ID：`grsai-wan-video`。
- Provider ID：`grsai-wan-video`。
- 能力：`video`。
- 默认 Base URL：`https://grsaiapi.com`。
- 鉴权驱动：`bearer`。
- 创建：`POST /v1/api/generate`。
- 查询：`GET /v1/api/result`。

## 配置字段

| 字段 | 类型 | 必填 | 含义 |
| --- | --- | --- | --- |
| `apiKey` | secret | 是 | API Key |

## 统一字段映射

| 统一字段 | 类型 | 必填 | 上游映射 | 说明 |
| --- | --- | --- | --- | --- |
| `model` | string | 是 | `model` | 仅 wan-3.0。 |
| `prompt` | string | 是 | `prompt` | 不可为空。 |
| `duration` | integer | 是 | `duration` | 生成时长 2–30 秒。 |
| `aspectRatio` | string | 是 | `aspectRatio` | 16:9、9:16、21:9、4:3、1:1、3:4 |
| `resolution` | string | 是 | `resolution` | 480p、720p、1080p |
| `images` | media[] | 否 | `images[] / firstFrame / lastFrame` | 参考图最多 10 张；显式首尾帧独立映射，不重复出现在 images 中。 |
| `videos` | media[] | 否 | `videos[].url` | 最多 5 段 MP4，总时长不超过 15 秒。 |
| `audios` | media[] | 否 | `audios[]` | 最多 5 段，每段最长 15 秒。 |
| `providerOptions` | object | 否 | `seed` | 仅接收本插件命名空间的整数 seed，包括 0。 |

## 上游请求模板逐字段清单

下表由插件请求模板生成，覆盖 body、query、headers 和 multipart 文件声明中的每个字段。

| 上游位置 | 值或转换表达式 |
| --- | --- |
| `create.method` | `"POST"` |
| `create.path` | `"/v1/api/generate"` |
| `create.contentType` | `"application/json"` |
| `create.body.model` | `{"$ref":"request.model"}` |
| `create.body.prompt` | `{"$ref":"request.prompt"}` |
| `create.body.aspectRatio` | `{"$ref":"request.output.aspectRatio"}` |
| `create.body.resolution` | `{"$ref":"request.output.resolution"}` |
| `create.body.duration` | `{"$ref":"request.output.duration"}` |
| `create.body.replyType` | `"async"` |
| `create.body.seed` | `{"$omitEmpty":{"$ref":"request.providerOptions.grsai-wan-video.seed"}}` |
| `create.body.firstFrame` | `{"$first":{"$map":{"from":{"$ref":"bindings.firstFrames"},"as":"media","in":{"$ref":"media.value"}}}}` |
| `create.body.lastFrame` | `{"$first":{"$map":{"from":{"$ref":"bindings.lastFrames"},"as":"media","in":{"$ref":"media.value"}}}}` |
| `create.body.images` | `{"$omitEmpty":{"$map":{"from":{"$ref":"bindings.referenceImages"},"as":"media","in":{"$ref":"media.value"}}}}` |
| `create.body.videos` | `{"$omitEmpty":{"$map":{"from":{"$ref":"bindings.videos"},"as":"media","in":{"url":{"$ref":"media.value"}}}}}` |
| `create.body.audios` | `{"$omitEmpty":{"$map":{"from":{"$ref":"bindings.audios"},"as":"media","in":{"$ref":"media.value"}}}}` |
| `create.originPath` | `true` |
| `poll.method` | `"GET"` |
| `poll.path` | `"/v1/api/result"` |
| `poll.originPath` | `true` |
| `poll.query.id` | `{"$ref":"taskId"}` |

## Provider 扩展键

- `providerOptions.grsai-wan-video.seed`

动态模型或工作流允许使用文档声明的完整 `parameters/input/extra_body` 对象；该对象是协议本身的开放 schema，不会被宿主裁剪。

## 响应映射逐字段清单

| 映射位置 | 上游路径或转换表达式 |
| --- | --- |
| `response.taskId` | `{"$if":{"condition":{"$and":[{"$eq":[{"$ref":"response.id"},{"$toString":{"$ref":"response.id"}}]},{"$gt":[{"$len":{"$trim":{"$ref":"response.id"}}},0]}]},"then":{"$ref":"response.id"},"else":null}}` |
| `response.status` | `{"$if":{"condition":{"$and":[{"$and":[{"$and":[{"$eq":[{"$ref":"response.id"},{"$toString":{"$ref":"response.id"}}]},{"$gt":[{"$len":{"$trim":{"$ref":"response.id"}}},0]}]},{"$or":[{"$not":{"$gt":[{"$len":{"$trim":{"$ref":"taskId"}}},0]}},{"$eq":[{"$ref":"taskId"},{"$ref":"response.id"}]}]}]},{"$eq":[{"$ref":"response.status"},"succeeded"]},{"$not":{"$gt":[{"$len":{"$trim":{"$ref":"response.error"}}},0]}},{"$and":[{"$eq":[{"$ref":"response.results"},{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item"}}}]},{"$gt":[{"$len":{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item.url"}}}},0]},{"$eq":[{"$len":{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item.url"}}}},{"$len":{"$ref":"response.results"}}]},{"$eq":[{"$len":{"$filter":{"from":{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item.url"}}},"as":"u","where":{"$not":{"$publicHttpsUrl":{"$ref":"u"}}}}}},0]}]}]},"then":"succeeded","else":{"$if":{"condition":{"$and":[{"$and":[{"$and":[{"$eq":[{"$ref":"response.id"},{"$toString":{"$ref":"response.id"}}]},{"$gt":[{"$len":{"$trim":{"$ref":"response.id"}}},0]}]},{"$or":[{"$not":{"$gt":[{"$len":{"$trim":{"$ref":"taskId"}}},0]}},{"$eq":[{"$ref":"taskId"},{"$ref":"response.id"}]}]}]},{"$eq":[{"$ref":"response.status"},"running"]}]},"then":"processing","else":"failed"}}}}` |
| `response.message` | `{"$coalesce":[{"$ref":"response.error"},{"$switch":{"cases":[{"when":{"$not":{"$and":[{"$and":[{"$eq":[{"$ref":"response.id"},{"$toString":{"$ref":"response.id"}}]},{"$gt":[{"$len":{"$trim":{"$ref":"response.id"}}},0]}]},{"$or":[{"$not":{"$gt":[{"$len":{"$trim":{"$ref":"taskId"}}},0]}},{"$eq":[{"$ref":"taskId"},{"$ref":"response.id"}]}]}]}},"then":"Grsai 返回了空或不匹配的任务 ID"},{"when":{"$eq":[{"$ref":"response.status"},"violation"]},"then":"Grsai 内容审核未通过"},{"when":{"$eq":[{"$ref":"response.status"},"failed"]},"then":"Grsai 视频生成失败"},{"when":{"$and":[{"$eq":[{"$ref":"response.status"},"succeeded"]},{"$not":{"$and":[{"$eq":[{"$ref":"response.results"},{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item"}}}]},{"$gt":[{"$len":{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item.url"}}}},0]},{"$eq":[{"$len":{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item.url"}}}},{"$len":{"$ref":"response.results"}}]},{"$eq":[{"$len":{"$filter":{"from":{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item.url"}}},"as":"u","where":{"$not":{"$publicHttpsUrl":{"$ref":"u"}}}}}},0]}]}}]},"then":"Grsai 完成响应缺少有效视频地址"},{"when":{"$not":{"$in":[{"$ref":"response.status"},["succeeded","running"]]}},"then":"Grsai 返回未知任务状态"}]}}]}` |
| `response.stage` | `{"$if":{"condition":{"$and":[{"$eq":[{"$ref":"response.status"},"running"]},{"$ne":[{"$ref":"response.progress"},null]},{"$eq":[{"$ref":"response.progress"},{"$toInt":{"$ref":"response.progress"}}]},{"$gte":[{"$ref":"response.progress"},0]},{"$lte":[{"$ref":"response.progress"},100]}]},"then":{"$concat":[{"$ref":"response.progress"},"%"]},"else":null}}` |
| `response.videos` | `{"$if":{"condition":{"$and":[{"$and":[{"$and":[{"$eq":[{"$ref":"response.id"},{"$toString":{"$ref":"response.id"}}]},{"$gt":[{"$len":{"$trim":{"$ref":"response.id"}}},0]}]},{"$or":[{"$not":{"$gt":[{"$len":{"$trim":{"$ref":"taskId"}}},0]}},{"$eq":[{"$ref":"taskId"},{"$ref":"response.id"}]}]}]},{"$eq":[{"$ref":"response.status"},"succeeded"]},{"$not":{"$gt":[{"$len":{"$trim":{"$ref":"response.error"}}},0]}},{"$and":[{"$eq":[{"$ref":"response.results"},{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item"}}}]},{"$gt":[{"$len":{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item.url"}}}},0]},{"$eq":[{"$len":{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item.url"}}}},{"$len":{"$ref":"response.results"}}]},{"$eq":[{"$len":{"$filter":{"from":{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item.url"}}},"as":"u","where":{"$not":{"$publicHttpsUrl":{"$ref":"u"}}}}}},0]}]}]},"then":{"$map":{"from":{"$ref":"response.results"},"as":"item","in":{"$ref":"item.url"}}},"else":null}}` |
| `response.usage` | `{"$ref":"response.usage"}` |
| `response.errorPaths[0]` | `"error"` |
| `response.resultEphemeral` | `true` |

## 响应与错误

插件把上游 task/status/text/media/usage 映射为统一结果。临时媒体 URL 标记为 ephemeral，由宿主立即下载持久化。HTTP 错误、业务 code 和 error object 保持失败语义，不包装成成功。

## 兼容边界

仅适配 [Grsai Wan](https://qmy27nhsd9.apifox.cn/531060669e0) 和[异步结果查询](https://qmy27nhsd9.apifox.cn/452409577e0)，并非 DashScope、RollDek 或 OpenAI 视频协议。全球节点为 https://grsaiapi.com，国内节点为 https://grsai.dakka.com.cn；渠道任选其一，同一任务在原节点查询。固定 replyType=async，不提供流式输出、回调、上游取消或模型列表。running 回执不是成功；成功读取全部 results[].url，违规和失败均终止。结果生成后仅 2 小时内能查询，不等于文件 URL 的有效期；宿主应持久化结果。可选 usage 原样保留，查询页未承诺返回 usage，不能把 credits 当 Token。供应商费用按（生成秒数＋参考视频秒数）×输出分辨率对应秒价；文档没有单价或小数取整规则。

普通参考图、音频使用字符串数组，视频使用 {url} 对象；素材按 order 排序。音频按单段最长 15 秒解释，视频按总和 15 秒；插件要求音视频带可信时长信息，拒绝无法检查的输入。宿主仅提交公网 HTTPS URL 或同类别的 Base64 data URI；不接收私网、HTTP 或无类型裸 Base64。上游明确首尾帧与 images、audios 互斥；示例混传不能作为有效组合。单独尾帧、首尾帧搭配参考视频未得到明确确认，本版主动收窄，不推断厂商禁用。未声明上游参考图编号语法。 纯音频参考模式表示“提示词＋参考音频”，不是无提示词；字段未要求必须附带画面素材，但纯音频组合需上游实测。视频编辑、增强等不在本接口范围；宿主的 video_to_video 价格选择器不是此处的请求模式。

<!-- YINGCE_MANIFEST_CONTRACT_START -->
## Manifest 完整接口定义

以下 JSON 与插件包内实际 `manifest.json` 逐字段一致，覆盖插件身份、权限、配置、鉴权、参数、校验、创建、Agent、查询、取消、结果下载、响应和 Agent 响应映射。`documentation` 字段的值就是当前完整文档；为避免文档在自身内部无限递归，JSON 中仅用等义占位文本表示正文。

```json
{
  "apiVersion": "yingce.plugin/v2",
  "id": "grsai-wan-video",
  "name": "Grsai Wan Video",
  "version": "1.0.0",
  "author": "影策社区",
  "description": "Grsai Wan Video 独立请求协议插件。",
  "documentation": "<当前插件的完整 documentation，由 README.md 与 docs/interface.md 拼接而成；为避免 JSON 递归，此处不重复展开正文。>",
  "permissions": [
    "generation.run",
    "media.read"
  ],
  "configuration": {
    "fields": [
      {
        "name": "apiKey",
        "type": "secret",
        "label": "API Key",
        "required": true
      }
    ]
  },
  "contributes": {
    "providers": [
      {
        "id": "grsai-wan-video",
        "label": "Grsai Wan Video",
        "capabilities": [
          "video"
        ],
        "scopes": [
          "admin.system-channel",
          "user.custom-channel",
          "canvas",
          "creation",
          "agent"
        ],
        "baseUrl": "https://grsaiapi.com",
        "requiresPublicMediaUrls": false,
        "detachedPolling": true,
        "auth": {
          "type": "bearer",
          "field": "apiKey"
        },
        "videoModes": {
          "wan-3.0": {
            "text_to_video": {
              "maxImages": 0,
              "inputRoles": [],
              "ratios": [
                "16:9",
                "9:16",
                "21:9",
                "4:3",
                "1:1",
                "3:4"
              ]
            },
            "image_to_video": {
              "maxImages": 2,
              "inputRoles": [
                "first_frame",
                "last_frame"
              ],
              "ratios": [
                "16:9",
                "9:16",
                "21:9",
                "4:3",
                "1:1",
                "3:4"
              ]
            },
            "reference_to_video": {
              "maxImages": 10,
              "inputRoles": [
                "reference_image",
                "reference_video",
                "reference_audio"
              ],
              "ratios": [
                "16:9",
                "9:16",
                "21:9",
                "4:3",
                "1:1",
                "3:4"
              ]
            },
            "audio_to_video": {
              "maxImages": 0,
              "inputRoles": [
                "reference_audio"
              ],
              "ratios": [
                "16:9",
                "9:16",
                "21:9",
                "4:3",
                "1:1",
                "3:4"
              ]
            }
          }
        },
        "videoOutputs": {
          "wan-3.0": {
            "resolutions": [
              "480p",
              "720p",
              "1080p"
            ],
            "minSeconds": 2,
            "maxSeconds": 30,
            "generateAudio": false
          }
        },
        "parameters": [
          {
            "name": "model",
            "type": "string",
            "required": true,
            "mapping": "model",
            "description": "仅 wan-3.0。"
          },
          {
            "name": "prompt",
            "type": "string",
            "required": true,
            "mapping": "prompt",
            "description": "不可为空。"
          },
          {
            "name": "duration",
            "type": "integer",
            "required": true,
            "mapping": "duration",
            "description": "生成时长 2–30 秒。"
          },
          {
            "name": "aspectRatio",
            "type": "string",
            "required": true,
            "mapping": "aspectRatio",
            "description": "16:9、9:16、21:9、4:3、1:1、3:4"
          },
          {
            "name": "resolution",
            "type": "string",
            "required": true,
            "mapping": "resolution",
            "description": "480p、720p、1080p"
          },
          {
            "name": "images",
            "type": "media[]",
            "required": false,
            "mapping": "images[] / firstFrame / lastFrame",
            "description": "参考图最多 10 张；显式首尾帧独立映射，不重复出现在 images 中。"
          },
          {
            "name": "videos",
            "type": "media[]",
            "required": false,
            "mapping": "videos[].url",
            "description": "最多 5 段 MP4，总时长不超过 15 秒。"
          },
          {
            "name": "audios",
            "type": "media[]",
            "required": false,
            "mapping": "audios[]",
            "description": "最多 5 段，每段最长 15 秒。"
          },
          {
            "name": "providerOptions",
            "type": "object",
            "required": false,
            "mapping": "seed",
            "description": "仅接收本插件命名空间的整数 seed，包括 0。"
          }
        ],
        "requestBindings": [
          {
            "name": "inputIDs",
            "value": {
              "$map": {
                "from": {
                  "$ref": "request.inputs"
                },
                "as": "media",
                "in": {
                  "kind": {
                    "$ref": "media.kind"
                  },
                  "role": {
                    "$ref": "media.role"
                  },
                  "id": {
                    "$ref": "media.id"
                  },
                  "source": {
                    "$if": {
                      "condition": {
                        "$not": {
                          "$gt": [
                            {
                              "$len": {
                                "$trim": {
                                  "$ref": "media.id"
                                }
                              }
                            },
                            0
                          ]
                        }
                      },
                      "then": {
                        "$ref": "media.value"
                      },
                      "else": null
                    }
                  }
                }
              }
            }
          },
          {
            "name": "typedIDs",
            "value": {
              "$map": {
                "from": {
                  "$concatArrays": [
                    {
                      "$ref": "request.images"
                    },
                    {
                      "$ref": "request.videos"
                    },
                    {
                      "$ref": "request.audios"
                    }
                  ]
                },
                "as": "media",
                "in": {
                  "kind": {
                    "$ref": "media.kind"
                  },
                  "role": {
                    "$ref": "media.role"
                  },
                  "id": {
                    "$ref": "media.id"
                  },
                  "source": {
                    "$if": {
                      "condition": {
                        "$not": {
                          "$gt": [
                            {
                              "$len": {
                                "$trim": {
                                  "$ref": "media.id"
                                }
                              }
                            },
                            0
                          ]
                        }
                      },
                      "then": {
                        "$ref": "media.value"
                      },
                      "else": null
                    }
                  }
                }
              }
            }
          },
          {
            "name": "images",
            "value": {
              "$sortByOrder": {
                "$if": {
                  "condition": {
                    "$gt": [
                      {
                        "$len": {
                          "$concatArrays": [
                            {
                              "$ref": "request.images"
                            },
                            {
                              "$ref": "request.videos"
                            },
                            {
                              "$ref": "request.audios"
                            }
                          ]
                        }
                      },
                      0
                    ]
                  },
                  "then": {
                    "$ref": "request.images"
                  },
                  "else": {
                    "$filter": {
                      "from": {
                        "$ref": "request.inputs"
                      },
                      "as": "media",
                      "where": {
                        "$eq": [
                          {
                            "$ref": "media.kind"
                          },
                          "image"
                        ]
                      }
                    }
                  }
                }
              }
            }
          },
          {
            "name": "videos",
            "value": {
              "$sortByOrder": {
                "$if": {
                  "condition": {
                    "$gt": [
                      {
                        "$len": {
                          "$concatArrays": [
                            {
                              "$ref": "request.images"
                            },
                            {
                              "$ref": "request.videos"
                            },
                            {
                              "$ref": "request.audios"
                            }
                          ]
                        }
                      },
                      0
                    ]
                  },
                  "then": {
                    "$ref": "request.videos"
                  },
                  "else": {
                    "$filter": {
                      "from": {
                        "$ref": "request.inputs"
                      },
                      "as": "media",
                      "where": {
                        "$eq": [
                          {
                            "$ref": "media.kind"
                          },
                          "video"
                        ]
                      }
                    }
                  }
                }
              }
            }
          },
          {
            "name": "audios",
            "value": {
              "$sortByOrder": {
                "$if": {
                  "condition": {
                    "$gt": [
                      {
                        "$len": {
                          "$concatArrays": [
                            {
                              "$ref": "request.images"
                            },
                            {
                              "$ref": "request.videos"
                            },
                            {
                              "$ref": "request.audios"
                            }
                          ]
                        }
                      },
                      0
                    ]
                  },
                  "then": {
                    "$ref": "request.audios"
                  },
                  "else": {
                    "$filter": {
                      "from": {
                        "$ref": "request.inputs"
                      },
                      "as": "media",
                      "where": {
                        "$eq": [
                          {
                            "$ref": "media.kind"
                          },
                          "audio"
                        ]
                      }
                    }
                  }
                }
              }
            }
          },
          {
            "name": "firstFrames",
            "value": {
              "$filter": {
                "from": {
                  "$ref": "bindings.images"
                },
                "as": "media",
                "where": {
                  "$in": [
                    {
                      "$ref": "media.role"
                    },
                    [
                      "first_frame"
                    ]
                  ]
                }
              }
            }
          },
          {
            "name": "lastFrames",
            "value": {
              "$filter": {
                "from": {
                  "$ref": "bindings.images"
                },
                "as": "media",
                "where": {
                  "$in": [
                    {
                      "$ref": "media.role"
                    },
                    [
                      "last_frame"
                    ]
                  ]
                }
              }
            }
          },
          {
            "name": "referenceImages",
            "value": {
              "$filter": {
                "from": {
                  "$ref": "bindings.images"
                },
                "as": "media",
                "where": {
                  "$in": [
                    {
                      "$ref": "media.role"
                    },
                    [
                      "",
                      "reference",
                      "reference_image",
                      "subject_reference",
                      "style_reference"
                    ]
                  ]
                }
              }
            }
          },
          {
            "name": "operation",
            "value": {
              "$coalesce": [
                {
                  "$ref": "request.operation"
                },
                {
                  "$if": {
                    "condition": {
                      "$gt": [
                        {
                          "$add": [
                            {
                              "$len": {
                                "$ref": "bindings.firstFrames"
                              }
                            },
                            {
                              "$len": {
                                "$ref": "bindings.lastFrames"
                              }
                            }
                          ]
                        },
                        0
                      ]
                    },
                    "then": "image_to_video",
                    "else": {
                      "$if": {
                        "condition": {
                          "$gt": [
                            {
                              "$add": [
                                {
                                  "$len": {
                                    "$ref": "bindings.images"
                                  }
                                },
                                {
                                  "$len": {
                                    "$ref": "bindings.videos"
                                  }
                                }
                              ]
                            },
                            0
                          ]
                        },
                        "then": "reference_to_video",
                        "else": {
                          "$if": {
                            "condition": {
                              "$gt": [
                                {
                                  "$len": {
                                    "$ref": "bindings.audios"
                                  }
                                },
                                0
                              ]
                            },
                            "then": "audio_to_video",
                            "else": "text_to_video"
                          }
                        }
                      }
                    }
                  }
                }
              ]
            }
          }
        ],
        "validations": [
          {
            "assert": {
              "$eq": [
                {
                  "$ref": "request.model"
                },
                "wan-3.0"
              ]
            },
            "message": "Grsai Wan · 仅支持 wan-3.0。"
          },
          {
            "assert": {
              "$gt": [
                {
                  "$len": {
                    "$trim": {
                      "$ref": "request.prompt"
                    }
                  }
                },
                0
              ]
            },
            "message": "Grsai Wan · 提示词不能为空。"
          },
          {
            "assert": {
              "$in": [
                {
                  "$ref": "request.output.aspectRatio"
                },
                [
                  "16:9",
                  "9:16",
                  "21:9",
                  "4:3",
                  "1:1",
                  "3:4"
                ]
              ]
            },
            "message": "Grsai Wan · 请选择支持的画幅比例。"
          },
          {
            "assert": {
              "$in": [
                {
                  "$ref": "request.output.resolution"
                },
                [
                  "480p",
                  "720p",
                  "1080p"
                ]
              ]
            },
            "message": "Grsai Wan · 分辨率须为 480p、720p 或 1080p。"
          },
          {
            "assert": {
              "$and": [
                {
                  "$gte": [
                    {
                      "$ref": "request.output.duration"
                    },
                    2
                  ]
                },
                {
                  "$lte": [
                    {
                      "$ref": "request.output.duration"
                    },
                    30
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 输出时长须为 2–30 秒。"
          },
          {
            "assert": {
              "$or": [
                {
                  "$eq": [
                    {
                      "$ref": "request.output.count"
                    },
                    null
                  ]
                },
                {
                  "$eq": [
                    {
                      "$ref": "request.output.count"
                    },
                    1
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 每个任务生成一个视频。"
          },
          {
            "assert": {
              "$and": [
                {
                  "$not": {
                    "$ref": "request.output.watermark"
                  }
                },
                {
                  "$not": {
                    "$ref": "request.output.generateAudio"
                  }
                }
              ]
            },
            "message": "Grsai Wan · 上游未提供水印或生成音轨开关。"
          },
          {
            "assert": {
              "$or": [
                {
                  "$eq": [
                    {
                      "$ref": "request.providerOptions.grsai-wan-video"
                    },
                    null
                  ]
                },
                {
                  "$eq": [
                    {
                      "$ref": "request.providerOptions.grsai-wan-video"
                    },
                    {
                      "seed": {
                        "$ref": "request.providerOptions.grsai-wan-video.seed"
                      }
                    }
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 不支持 seed 以外的扩展参数或空值。"
          },
          {
            "assert": {
              "$or": [
                {
                  "$eq": [
                    {
                      "$ref": "request.providerOptions.grsai-wan-video.seed"
                    },
                    null
                  ]
                },
                {
                  "$eq": [
                    {
                      "$ref": "request.providerOptions.grsai-wan-video.seed"
                    },
                    {
                      "$toInt": {
                        "$ref": "request.providerOptions.grsai-wan-video.seed"
                      }
                    }
                  ]
                }
              ]
            },
            "message": "Grsai Wan · seed 必须为整数。"
          },
          {
            "assert": {
              "$in": [
                {
                  "$ref": "bindings.operation"
                },
                [
                  "text_to_video",
                  "image_to_video",
                  "reference_to_video",
                  "audio_to_video"
                ]
              ]
            },
            "message": "Grsai Wan · 仅支持文生、关键帧、全能参考或音频参考生成模式。"
          },
          {
            "assert": {
              "$lte": [
                {
                  "$len": {
                    "$ref": "request.inputs"
                  }
                },
                20
              ]
            },
            "message": "Grsai Wan · 参考素材数量超过上限。"
          },
          {
            "assert": {
              "$if": {
                "condition": {
                  "$eq": [
                    {
                      "$len": {
                        "$ref": "bindings.typedIDs"
                      }
                    },
                    0
                  ]
                },
                "then": true,
                "else": {
                  "$if": {
                    "condition": {
                      "$eq": [
                        {
                          "$ref": "bindings.inputIDs"
                        },
                        {
                          "$ref": "bindings.typedIDs"
                        }
                      ]
                    },
                    "then": true,
                    "else": {
                      "$and": [
                        {
                          "$eq": [
                            {
                              "$len": {
                                "$ref": "bindings.inputIDs"
                              }
                            },
                            {
                              "$len": {
                                "$ref": "bindings.typedIDs"
                              }
                            }
                          ]
                        },
                        {
                          "$eq": [
                            {
                              "$len": {
                                "$filter": {
                                  "from": {
                                    "$ref": "bindings.inputIDs"
                                  },
                                  "as": "candidate",
                                  "where": {
                                    "$ne": [
                                      {
                                        "$len": {
                                          "$filter": {
                                            "from": {
                                              "$ref": "bindings.inputIDs"
                                            },
                                            "as": "other",
                                            "where": {
                                              "$eq": [
                                                {
                                                  "$ref": "candidate"
                                                },
                                                {
                                                  "$ref": "other"
                                                }
                                              ]
                                            }
                                          }
                                        }
                                      },
                                      {
                                        "$len": {
                                          "$filter": {
                                            "from": {
                                              "$ref": "bindings.typedIDs"
                                            },
                                            "as": "other",
                                            "where": {
                                              "$eq": [
                                                {
                                                  "$ref": "candidate"
                                                },
                                                {
                                                  "$ref": "other"
                                                }
                                              ]
                                            }
                                          }
                                        }
                                      }
                                    ]
                                  }
                                }
                              }
                            },
                            0
                          ]
                        }
                      ]
                    }
                  }
                }
              }
            },
            "message": "Grsai Wan · 两份素材列表不一致。"
          },
          {
            "assert": {
              "$eq": [
                {
                  "$len": {
                    "$ref": "request.inputs"
                  }
                },
                {
                  "$add": [
                    {
                      "$len": {
                        "$ref": "bindings.images"
                      }
                    },
                    {
                      "$len": {
                        "$ref": "bindings.videos"
                      }
                    },
                    {
                      "$len": {
                        "$ref": "bindings.audios"
                      }
                    }
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 仅支持图片、视频与音频素材。"
          },
          {
            "assert": {
              "$eq": [
                {
                  "$len": {
                    "$filter": {
                      "from": {
                        "$ref": "bindings.images"
                      },
                      "as": "media",
                      "where": {
                        "$not": {
                          "$and": [
                            {
                              "$in": [
                                {
                                  "$ref": "media.kind"
                                },
                                [
                                  "",
                                  "image"
                                ]
                              ]
                            },
                            {
                              "$in": [
                                {
                                  "$ref": "media.role"
                                },
                                [
                                  "",
                                  "reference",
                                  "reference_image",
                                  "subject_reference",
                                  "style_reference",
                                  "first_frame",
                                  "last_frame"
                                ]
                              ]
                            }
                          ]
                        }
                      }
                    }
                  }
                },
                0
              ]
            },
            "message": "Grsai Wan · 图片角色无效。"
          },
          {
            "assert": {
              "$eq": [
                {
                  "$len": {
                    "$filter": {
                      "from": {
                        "$ref": "bindings.videos"
                      },
                      "as": "media",
                      "where": {
                        "$not": {
                          "$and": [
                            {
                              "$in": [
                                {
                                  "$ref": "media.kind"
                                },
                                [
                                  "",
                                  "video"
                                ]
                              ]
                            },
                            {
                              "$in": [
                                {
                                  "$ref": "media.role"
                                },
                                [
                                  "",
                                  "reference",
                                  "reference_video"
                                ]
                              ]
                            }
                          ]
                        }
                      }
                    }
                  }
                },
                0
              ]
            },
            "message": "Grsai Wan · 视频角色无效。"
          },
          {
            "assert": {
              "$eq": [
                {
                  "$len": {
                    "$filter": {
                      "from": {
                        "$ref": "bindings.audios"
                      },
                      "as": "media",
                      "where": {
                        "$not": {
                          "$and": [
                            {
                              "$in": [
                                {
                                  "$ref": "media.kind"
                                },
                                [
                                  "",
                                  "audio"
                                ]
                              ]
                            },
                            {
                              "$in": [
                                {
                                  "$ref": "media.role"
                                },
                                [
                                  "",
                                  "reference",
                                  "reference_audio",
                                  "reference_voice"
                                ]
                              ]
                            }
                          ]
                        }
                      }
                    }
                  }
                },
                0
              ]
            },
            "message": "Grsai Wan · 音频角色无效。"
          },
          {
            "assert": {
              "$and": [
                {
                  "$lte": [
                    {
                      "$len": {
                        "$ref": "bindings.firstFrames"
                      }
                    },
                    1
                  ]
                },
                {
                  "$lte": [
                    {
                      "$len": {
                        "$ref": "bindings.lastFrames"
                      }
                    },
                    1
                  ]
                },
                {
                  "$lte": [
                    {
                      "$len": {
                        "$ref": "bindings.referenceImages"
                      }
                    },
                    10
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 首尾帧各最多 1 张，普通参考图最多 10 张。"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$gt": [
                      {
                        "$add": [
                          {
                            "$len": {
                              "$ref": "bindings.firstFrames"
                            }
                          },
                          {
                            "$len": {
                              "$ref": "bindings.lastFrames"
                            }
                          }
                        ]
                      },
                      0
                    ]
                  }
                },
                {
                  "$eq": [
                    {
                      "$add": [
                        {
                          "$len": {
                            "$ref": "bindings.referenceImages"
                          }
                        },
                        {
                          "$len": {
                            "$ref": "bindings.audios"
                          }
                        }
                      ]
                    },
                    0
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 首尾帧与参考图、参考音频互斥。"
          },
          {
            "assert": {
              "$or": [
                {
                  "$not": {
                    "$gt": [
                      {
                        "$add": [
                          {
                            "$len": {
                              "$ref": "bindings.firstFrames"
                            }
                          },
                          {
                            "$len": {
                              "$ref": "bindings.lastFrames"
                            }
                          }
                        ]
                      },
                      0
                    ]
                  }
                },
                {
                  "$and": [
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "bindings.firstFrames"
                          }
                        },
                        1
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$len": {
                            "$ref": "bindings.videos"
                          }
                        },
                        0
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 本插件暂不开放单独尾帧或首尾帧混合参考视频。"
          },
          {
            "assert": {
              "$and": [
                {
                  "$lte": [
                    {
                      "$len": {
                        "$ref": "bindings.videos"
                      }
                    },
                    5
                  ]
                },
                {
                  "$lte": [
                    {
                      "$len": {
                        "$ref": "bindings.audios"
                      }
                    },
                    5
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 参考视频和音频分别最多 5 段。"
          },
          {
            "assert": {
              "$eq": [
                {
                  "$len": {
                    "$filter": {
                      "from": {
                        "$ref": "bindings.images"
                      },
                      "as": "media",
                      "where": {
                        "$not": {
                          "$and": [
                            {
                              "$or": [
                                {
                                  "$eq": [
                                    {
                                      "$lower": {
                                        "$trim": {
                                          "$ref": "media.mimeType"
                                        }
                                      }
                                    },
                                    ""
                                  ]
                                },
                                {
                                  "$eq": [
                                    {
                                      "$lower": {
                                        "$trim": {
                                          "$ref": "media.mimeType"
                                        }
                                      }
                                    },
                                    "image"
                                  ]
                                },
                                {
                                  "$eq": [
                                    {
                                      "$first": {
                                        "$split": [
                                          {
                                            "$lower": {
                                              "$trim": {
                                                "$ref": "media.mimeType"
                                              }
                                            }
                                          },
                                          "/"
                                        ]
                                      }
                                    },
                                    "image"
                                  ]
                                }
                              ]
                            },
                            {
                              "$if": {
                                "condition": {
                                  "$publicHttpsUrl": {
                                    "$ref": "media.value"
                                  }
                                },
                                "then": true,
                                "else": {
                                  "$if": {
                                    "condition": {
                                      "$eq": [
                                        {
                                          "$len": {
                                            "$split": [
                                              {
                                                "$ref": "media.value"
                                              },
                                              ","
                                            ]
                                          }
                                        },
                                        2
                                      ]
                                    },
                                    "then": {
                                      "$and": [
                                        {
                                          "$eq": [
                                            {
                                              "$countMatches": [
                                                {
                                                  "$first": {
                                                    "$split": [
                                                      {
                                                        "$ref": "media.value"
                                                      },
                                                      ","
                                                    ]
                                                  }
                                                },
                                                "^data:image/[A-Za-z0-9.+-]+;base64$"
                                              ]
                                            },
                                            1
                                          ]
                                        },
                                        {
                                          "$gt": [
                                            {
                                              "$len": {
                                                "$dataPayload": {
                                                  "$ref": "media.value"
                                                }
                                              }
                                            },
                                            0
                                          ]
                                        }
                                      ]
                                    },
                                    "else": false
                                  }
                                }
                              }
                            }
                          ]
                        }
                      }
                    }
                  }
                },
                0
              ]
            },
            "message": "Grsai Wan · 图片须使用公网 HTTPS URL 或图片 Base64 data URI。"
          },
          {
            "assert": {
              "$eq": [
                {
                  "$len": {
                    "$filter": {
                      "from": {
                        "$ref": "bindings.audios"
                      },
                      "as": "media",
                      "where": {
                        "$not": {
                          "$and": [
                            {
                              "$or": [
                                {
                                  "$eq": [
                                    {
                                      "$lower": {
                                        "$trim": {
                                          "$ref": "media.mimeType"
                                        }
                                      }
                                    },
                                    ""
                                  ]
                                },
                                {
                                  "$eq": [
                                    {
                                      "$lower": {
                                        "$trim": {
                                          "$ref": "media.mimeType"
                                        }
                                      }
                                    },
                                    "audio"
                                  ]
                                },
                                {
                                  "$eq": [
                                    {
                                      "$first": {
                                        "$split": [
                                          {
                                            "$lower": {
                                              "$trim": {
                                                "$ref": "media.mimeType"
                                              }
                                            }
                                          },
                                          "/"
                                        ]
                                      }
                                    },
                                    "audio"
                                  ]
                                }
                              ]
                            },
                            {
                              "$if": {
                                "condition": {
                                  "$publicHttpsUrl": {
                                    "$ref": "media.value"
                                  }
                                },
                                "then": true,
                                "else": {
                                  "$if": {
                                    "condition": {
                                      "$eq": [
                                        {
                                          "$len": {
                                            "$split": [
                                              {
                                                "$ref": "media.value"
                                              },
                                              ","
                                            ]
                                          }
                                        },
                                        2
                                      ]
                                    },
                                    "then": {
                                      "$and": [
                                        {
                                          "$eq": [
                                            {
                                              "$countMatches": [
                                                {
                                                  "$first": {
                                                    "$split": [
                                                      {
                                                        "$ref": "media.value"
                                                      },
                                                      ","
                                                    ]
                                                  }
                                                },
                                                "^data:audio/[A-Za-z0-9.+-]+;base64$"
                                              ]
                                            },
                                            1
                                          ]
                                        },
                                        {
                                          "$gt": [
                                            {
                                              "$len": {
                                                "$dataPayload": {
                                                  "$ref": "media.value"
                                                }
                                              }
                                            },
                                            0
                                          ]
                                        }
                                      ]
                                    },
                                    "else": false
                                  }
                                }
                              }
                            }
                          ]
                        }
                      }
                    }
                  }
                },
                0
              ]
            },
            "message": "Grsai Wan · 音频须使用公网 HTTPS URL 或音频 Base64 data URI。"
          },
          {
            "assert": {
              "$eq": [
                {
                  "$len": {
                    "$filter": {
                      "from": {
                        "$ref": "bindings.videos"
                      },
                      "as": "media",
                      "where": {
                        "$not": {
                          "$and": [
                            {
                              "$or": [
                                {
                                  "$eq": [
                                    {
                                      "$lower": {
                                        "$trim": {
                                          "$ref": "media.mimeType"
                                        }
                                      }
                                    },
                                    ""
                                  ]
                                },
                                {
                                  "$eq": [
                                    {
                                      "$lower": {
                                        "$trim": {
                                          "$ref": "media.mimeType"
                                        }
                                      }
                                    },
                                    "video/mp4"
                                  ]
                                }
                              ]
                            },
                            {
                              "$if": {
                                "condition": {
                                  "$publicHttpsUrl": {
                                    "$ref": "media.value"
                                  }
                                },
                                "then": {
                                  "$or": [
                                    {
                                      "$eq": [
                                        {
                                          "$lower": {
                                            "$trim": {
                                              "$ref": "media.mimeType"
                                            }
                                          }
                                        },
                                        "video/mp4"
                                      ]
                                    },
                                    {
                                      "$countMatches": [
                                        {
                                          "$ref": "media.value"
                                        },
                                        "(?i)^https://[^?#]+\\.mp4(?:[?#].*)?$"
                                      ]
                                    }
                                  ]
                                },
                                "else": {
                                  "$if": {
                                    "condition": {
                                      "$eq": [
                                        {
                                          "$len": {
                                            "$split": [
                                              {
                                                "$ref": "media.value"
                                              },
                                              ","
                                            ]
                                          }
                                        },
                                        2
                                      ]
                                    },
                                    "then": {
                                      "$and": [
                                        {
                                          "$eq": [
                                            {
                                              "$countMatches": [
                                                {
                                                  "$first": {
                                                    "$split": [
                                                      {
                                                        "$ref": "media.value"
                                                      },
                                                      ","
                                                    ]
                                                  }
                                                },
                                                "^data:video/mp4;base64$"
                                              ]
                                            },
                                            1
                                          ]
                                        },
                                        {
                                          "$gt": [
                                            {
                                              "$len": {
                                                "$dataPayload": {
                                                  "$ref": "media.value"
                                                }
                                              }
                                            },
                                            0
                                          ]
                                        }
                                      ]
                                    },
                                    "else": false
                                  }
                                }
                              }
                            }
                          ]
                        }
                      }
                    }
                  }
                },
                0
              ]
            },
            "message": "Grsai Wan · 参考视频仅支持 MP4；请提供 video/mp4 或可识别的 MP4 地址。"
          },
          {
            "assert": {
              "$and": [
                {
                  "$eq": [
                    {
                      "$len": {
                        "$filter": {
                          "from": {
                            "$ref": "bindings.videos"
                          },
                          "as": "media",
                          "where": {
                            "$not": {
                              "$gt": [
                                {
                                  "$ref": "media.metadata.durationMs"
                                },
                                0
                              ]
                            }
                          }
                        }
                      }
                    },
                    0
                  ]
                },
                {
                  "$lte": [
                    {
                      "$sum": {
                        "$map": {
                          "from": {
                            "$ref": "bindings.videos"
                          },
                          "as": "media",
                          "in": {
                            "$ref": "media.metadata.durationMs"
                          }
                        }
                      }
                    },
                    15000
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 参考视频须有时长，总时长不超过 15 秒。"
          },
          {
            "assert": {
              "$eq": [
                {
                  "$len": {
                    "$filter": {
                      "from": {
                        "$ref": "bindings.audios"
                      },
                      "as": "media",
                      "where": {
                        "$not": {
                          "$and": [
                            {
                              "$gt": [
                                {
                                  "$ref": "media.metadata.durationMs"
                                },
                                0
                              ]
                            },
                            {
                              "$lte": [
                                {
                                  "$ref": "media.metadata.durationMs"
                                },
                                15000
                              ]
                            }
                          ]
                        }
                      }
                    }
                  }
                },
                0
              ]
            },
            "message": "Grsai Wan · 参考音频须有时长，每段不超过 15 秒。"
          },
          {
            "assert": {
              "$or": [
                {
                  "$ne": [
                    {
                      "$ref": "bindings.operation"
                    },
                    "text_to_video"
                  ]
                },
                {
                  "$eq": [
                    {
                      "$len": {
                        "$ref": "request.inputs"
                      }
                    },
                    0
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 文生视频不能携带素材。"
          },
          {
            "assert": {
              "$or": [
                {
                  "$ne": [
                    {
                      "$ref": "bindings.operation"
                    },
                    "image_to_video"
                  ]
                },
                {
                  "$gt": [
                    {
                      "$add": [
                        {
                          "$len": {
                            "$ref": "bindings.firstFrames"
                          }
                        },
                        {
                          "$len": {
                            "$ref": "bindings.lastFrames"
                          }
                        }
                      ]
                    },
                    0
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 图生视频须提供首帧。"
          },
          {
            "assert": {
              "$or": [
                {
                  "$ne": [
                    {
                      "$ref": "bindings.operation"
                    },
                    "reference_to_video"
                  ]
                },
                {
                  "$and": [
                    {
                      "$not": {
                        "$gt": [
                          {
                            "$add": [
                              {
                                "$len": {
                                  "$ref": "bindings.firstFrames"
                                }
                              },
                              {
                                "$len": {
                                  "$ref": "bindings.lastFrames"
                                }
                              }
                            ]
                          },
                          0
                        ]
                      }
                    },
                    {
                      "$gt": [
                        {
                          "$len": {
                            "$ref": "request.inputs"
                          }
                        },
                        0
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 全能参考须提供普通参考素材，不能使用首尾帧。"
          },
          {
            "assert": {
              "$or": [
                {
                  "$ne": [
                    {
                      "$ref": "bindings.operation"
                    },
                    "audio_to_video"
                  ]
                },
                {
                  "$and": [
                    {
                      "$gt": [
                        {
                          "$len": {
                            "$ref": "bindings.audios"
                          }
                        },
                        0
                      ]
                    },
                    {
                      "$eq": [
                        {
                          "$add": [
                            {
                              "$len": {
                                "$ref": "bindings.images"
                              }
                            },
                            {
                              "$len": {
                                "$ref": "bindings.videos"
                              }
                            }
                          ]
                        },
                        0
                      ]
                    }
                  ]
                }
              ]
            },
            "message": "Grsai Wan · 音频参考模式只接受参考音频。"
          }
        ],
        "create": {
          "method": "POST",
          "path": "/v1/api/generate",
          "contentType": "application/json",
          "body": {
            "model": {
              "$ref": "request.model"
            },
            "prompt": {
              "$ref": "request.prompt"
            },
            "aspectRatio": {
              "$ref": "request.output.aspectRatio"
            },
            "resolution": {
              "$ref": "request.output.resolution"
            },
            "duration": {
              "$ref": "request.output.duration"
            },
            "replyType": "async",
            "seed": {
              "$omitEmpty": {
                "$ref": "request.providerOptions.grsai-wan-video.seed"
              }
            },
            "firstFrame": {
              "$first": {
                "$map": {
                  "from": {
                    "$ref": "bindings.firstFrames"
                  },
                  "as": "media",
                  "in": {
                    "$ref": "media.value"
                  }
                }
              }
            },
            "lastFrame": {
              "$first": {
                "$map": {
                  "from": {
                    "$ref": "bindings.lastFrames"
                  },
                  "as": "media",
                  "in": {
                    "$ref": "media.value"
                  }
                }
              }
            },
            "images": {
              "$omitEmpty": {
                "$map": {
                  "from": {
                    "$ref": "bindings.referenceImages"
                  },
                  "as": "media",
                  "in": {
                    "$ref": "media.value"
                  }
                }
              }
            },
            "videos": {
              "$omitEmpty": {
                "$map": {
                  "from": {
                    "$ref": "bindings.videos"
                  },
                  "as": "media",
                  "in": {
                    "url": {
                      "$ref": "media.value"
                    }
                  }
                }
              }
            },
            "audios": {
              "$omitEmpty": {
                "$map": {
                  "from": {
                    "$ref": "bindings.audios"
                  },
                  "as": "media",
                  "in": {
                    "$ref": "media.value"
                  }
                }
              }
            }
          },
          "originPath": true
        },
        "poll": {
          "method": "GET",
          "path": "/v1/api/result",
          "originPath": true,
          "query": {
            "id": {
              "$ref": "taskId"
            }
          }
        },
        "response": {
          "taskId": {
            "$if": {
              "condition": {
                "$and": [
                  {
                    "$eq": [
                      {
                        "$ref": "response.id"
                      },
                      {
                        "$toString": {
                          "$ref": "response.id"
                        }
                      }
                    ]
                  },
                  {
                    "$gt": [
                      {
                        "$len": {
                          "$trim": {
                            "$ref": "response.id"
                          }
                        }
                      },
                      0
                    ]
                  }
                ]
              },
              "then": {
                "$ref": "response.id"
              },
              "else": null
            }
          },
          "status": {
            "$if": {
              "condition": {
                "$and": [
                  {
                    "$and": [
                      {
                        "$and": [
                          {
                            "$eq": [
                              {
                                "$ref": "response.id"
                              },
                              {
                                "$toString": {
                                  "$ref": "response.id"
                                }
                              }
                            ]
                          },
                          {
                            "$gt": [
                              {
                                "$len": {
                                  "$trim": {
                                    "$ref": "response.id"
                                  }
                                }
                              },
                              0
                            ]
                          }
                        ]
                      },
                      {
                        "$or": [
                          {
                            "$not": {
                              "$gt": [
                                {
                                  "$len": {
                                    "$trim": {
                                      "$ref": "taskId"
                                    }
                                  }
                                },
                                0
                              ]
                            }
                          },
                          {
                            "$eq": [
                              {
                                "$ref": "taskId"
                              },
                              {
                                "$ref": "response.id"
                              }
                            ]
                          }
                        ]
                      }
                    ]
                  },
                  {
                    "$eq": [
                      {
                        "$ref": "response.status"
                      },
                      "succeeded"
                    ]
                  },
                  {
                    "$not": {
                      "$gt": [
                        {
                          "$len": {
                            "$trim": {
                              "$ref": "response.error"
                            }
                          }
                        },
                        0
                      ]
                    }
                  },
                  {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "response.results"
                          },
                          {
                            "$map": {
                              "from": {
                                "$ref": "response.results"
                              },
                              "as": "item",
                              "in": {
                                "$ref": "item"
                              }
                            }
                          }
                        ]
                      },
                      {
                        "$gt": [
                          {
                            "$len": {
                              "$map": {
                                "from": {
                                  "$ref": "response.results"
                                },
                                "as": "item",
                                "in": {
                                  "$ref": "item.url"
                                }
                              }
                            }
                          },
                          0
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$len": {
                              "$map": {
                                "from": {
                                  "$ref": "response.results"
                                },
                                "as": "item",
                                "in": {
                                  "$ref": "item.url"
                                }
                              }
                            }
                          },
                          {
                            "$len": {
                              "$ref": "response.results"
                            }
                          }
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$len": {
                              "$filter": {
                                "from": {
                                  "$map": {
                                    "from": {
                                      "$ref": "response.results"
                                    },
                                    "as": "item",
                                    "in": {
                                      "$ref": "item.url"
                                    }
                                  }
                                },
                                "as": "u",
                                "where": {
                                  "$not": {
                                    "$publicHttpsUrl": {
                                      "$ref": "u"
                                    }
                                  }
                                }
                              }
                            }
                          },
                          0
                        ]
                      }
                    ]
                  }
                ]
              },
              "then": "succeeded",
              "else": {
                "$if": {
                  "condition": {
                    "$and": [
                      {
                        "$and": [
                          {
                            "$and": [
                              {
                                "$eq": [
                                  {
                                    "$ref": "response.id"
                                  },
                                  {
                                    "$toString": {
                                      "$ref": "response.id"
                                    }
                                  }
                                ]
                              },
                              {
                                "$gt": [
                                  {
                                    "$len": {
                                      "$trim": {
                                        "$ref": "response.id"
                                      }
                                    }
                                  },
                                  0
                                ]
                              }
                            ]
                          },
                          {
                            "$or": [
                              {
                                "$not": {
                                  "$gt": [
                                    {
                                      "$len": {
                                        "$trim": {
                                          "$ref": "taskId"
                                        }
                                      }
                                    },
                                    0
                                  ]
                                }
                              },
                              {
                                "$eq": [
                                  {
                                    "$ref": "taskId"
                                  },
                                  {
                                    "$ref": "response.id"
                                  }
                                ]
                              }
                            ]
                          }
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$ref": "response.status"
                          },
                          "running"
                        ]
                      }
                    ]
                  },
                  "then": "processing",
                  "else": "failed"
                }
              }
            }
          },
          "message": {
            "$coalesce": [
              {
                "$ref": "response.error"
              },
              {
                "$switch": {
                  "cases": [
                    {
                      "when": {
                        "$not": {
                          "$and": [
                            {
                              "$and": [
                                {
                                  "$eq": [
                                    {
                                      "$ref": "response.id"
                                    },
                                    {
                                      "$toString": {
                                        "$ref": "response.id"
                                      }
                                    }
                                  ]
                                },
                                {
                                  "$gt": [
                                    {
                                      "$len": {
                                        "$trim": {
                                          "$ref": "response.id"
                                        }
                                      }
                                    },
                                    0
                                  ]
                                }
                              ]
                            },
                            {
                              "$or": [
                                {
                                  "$not": {
                                    "$gt": [
                                      {
                                        "$len": {
                                          "$trim": {
                                            "$ref": "taskId"
                                          }
                                        }
                                      },
                                      0
                                    ]
                                  }
                                },
                                {
                                  "$eq": [
                                    {
                                      "$ref": "taskId"
                                    },
                                    {
                                      "$ref": "response.id"
                                    }
                                  ]
                                }
                              ]
                            }
                          ]
                        }
                      },
                      "then": "Grsai 返回了空或不匹配的任务 ID"
                    },
                    {
                      "when": {
                        "$eq": [
                          {
                            "$ref": "response.status"
                          },
                          "violation"
                        ]
                      },
                      "then": "Grsai 内容审核未通过"
                    },
                    {
                      "when": {
                        "$eq": [
                          {
                            "$ref": "response.status"
                          },
                          "failed"
                        ]
                      },
                      "then": "Grsai 视频生成失败"
                    },
                    {
                      "when": {
                        "$and": [
                          {
                            "$eq": [
                              {
                                "$ref": "response.status"
                              },
                              "succeeded"
                            ]
                          },
                          {
                            "$not": {
                              "$and": [
                                {
                                  "$eq": [
                                    {
                                      "$ref": "response.results"
                                    },
                                    {
                                      "$map": {
                                        "from": {
                                          "$ref": "response.results"
                                        },
                                        "as": "item",
                                        "in": {
                                          "$ref": "item"
                                        }
                                      }
                                    }
                                  ]
                                },
                                {
                                  "$gt": [
                                    {
                                      "$len": {
                                        "$map": {
                                          "from": {
                                            "$ref": "response.results"
                                          },
                                          "as": "item",
                                          "in": {
                                            "$ref": "item.url"
                                          }
                                        }
                                      }
                                    },
                                    0
                                  ]
                                },
                                {
                                  "$eq": [
                                    {
                                      "$len": {
                                        "$map": {
                                          "from": {
                                            "$ref": "response.results"
                                          },
                                          "as": "item",
                                          "in": {
                                            "$ref": "item.url"
                                          }
                                        }
                                      }
                                    },
                                    {
                                      "$len": {
                                        "$ref": "response.results"
                                      }
                                    }
                                  ]
                                },
                                {
                                  "$eq": [
                                    {
                                      "$len": {
                                        "$filter": {
                                          "from": {
                                            "$map": {
                                              "from": {
                                                "$ref": "response.results"
                                              },
                                              "as": "item",
                                              "in": {
                                                "$ref": "item.url"
                                              }
                                            }
                                          },
                                          "as": "u",
                                          "where": {
                                            "$not": {
                                              "$publicHttpsUrl": {
                                                "$ref": "u"
                                              }
                                            }
                                          }
                                        }
                                      }
                                    },
                                    0
                                  ]
                                }
                              ]
                            }
                          }
                        ]
                      },
                      "then": "Grsai 完成响应缺少有效视频地址"
                    },
                    {
                      "when": {
                        "$not": {
                          "$in": [
                            {
                              "$ref": "response.status"
                            },
                            [
                              "succeeded",
                              "running"
                            ]
                          ]
                        }
                      },
                      "then": "Grsai 返回未知任务状态"
                    }
                  ]
                }
              }
            ]
          },
          "stage": {
            "$if": {
              "condition": {
                "$and": [
                  {
                    "$eq": [
                      {
                        "$ref": "response.status"
                      },
                      "running"
                    ]
                  },
                  {
                    "$ne": [
                      {
                        "$ref": "response.progress"
                      },
                      null
                    ]
                  },
                  {
                    "$eq": [
                      {
                        "$ref": "response.progress"
                      },
                      {
                        "$toInt": {
                          "$ref": "response.progress"
                        }
                      }
                    ]
                  },
                  {
                    "$gte": [
                      {
                        "$ref": "response.progress"
                      },
                      0
                    ]
                  },
                  {
                    "$lte": [
                      {
                        "$ref": "response.progress"
                      },
                      100
                    ]
                  }
                ]
              },
              "then": {
                "$concat": [
                  {
                    "$ref": "response.progress"
                  },
                  "%"
                ]
              },
              "else": null
            }
          },
          "videos": {
            "$if": {
              "condition": {
                "$and": [
                  {
                    "$and": [
                      {
                        "$and": [
                          {
                            "$eq": [
                              {
                                "$ref": "response.id"
                              },
                              {
                                "$toString": {
                                  "$ref": "response.id"
                                }
                              }
                            ]
                          },
                          {
                            "$gt": [
                              {
                                "$len": {
                                  "$trim": {
                                    "$ref": "response.id"
                                  }
                                }
                              },
                              0
                            ]
                          }
                        ]
                      },
                      {
                        "$or": [
                          {
                            "$not": {
                              "$gt": [
                                {
                                  "$len": {
                                    "$trim": {
                                      "$ref": "taskId"
                                    }
                                  }
                                },
                                0
                              ]
                            }
                          },
                          {
                            "$eq": [
                              {
                                "$ref": "taskId"
                              },
                              {
                                "$ref": "response.id"
                              }
                            ]
                          }
                        ]
                      }
                    ]
                  },
                  {
                    "$eq": [
                      {
                        "$ref": "response.status"
                      },
                      "succeeded"
                    ]
                  },
                  {
                    "$not": {
                      "$gt": [
                        {
                          "$len": {
                            "$trim": {
                              "$ref": "response.error"
                            }
                          }
                        },
                        0
                      ]
                    }
                  },
                  {
                    "$and": [
                      {
                        "$eq": [
                          {
                            "$ref": "response.results"
                          },
                          {
                            "$map": {
                              "from": {
                                "$ref": "response.results"
                              },
                              "as": "item",
                              "in": {
                                "$ref": "item"
                              }
                            }
                          }
                        ]
                      },
                      {
                        "$gt": [
                          {
                            "$len": {
                              "$map": {
                                "from": {
                                  "$ref": "response.results"
                                },
                                "as": "item",
                                "in": {
                                  "$ref": "item.url"
                                }
                              }
                            }
                          },
                          0
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$len": {
                              "$map": {
                                "from": {
                                  "$ref": "response.results"
                                },
                                "as": "item",
                                "in": {
                                  "$ref": "item.url"
                                }
                              }
                            }
                          },
                          {
                            "$len": {
                              "$ref": "response.results"
                            }
                          }
                        ]
                      },
                      {
                        "$eq": [
                          {
                            "$len": {
                              "$filter": {
                                "from": {
                                  "$map": {
                                    "from": {
                                      "$ref": "response.results"
                                    },
                                    "as": "item",
                                    "in": {
                                      "$ref": "item.url"
                                    }
                                  }
                                },
                                "as": "u",
                                "where": {
                                  "$not": {
                                    "$publicHttpsUrl": {
                                      "$ref": "u"
                                    }
                                  }
                                }
                              }
                            }
                          },
                          0
                        ]
                      }
                    ]
                  }
                ]
              },
              "then": {
                "$map": {
                  "from": {
                    "$ref": "response.results"
                  },
                  "as": "item",
                  "in": {
                    "$ref": "item.url"
                  }
                }
              },
              "else": null
            }
          },
          "usage": {
            "$ref": "response.usage"
          },
          "errorPaths": [
            "error"
          ],
          "resultEphemeral": true
        }
      }
    ]
  }
}
```
<!-- YINGCE_MANIFEST_CONTRACT_END -->
