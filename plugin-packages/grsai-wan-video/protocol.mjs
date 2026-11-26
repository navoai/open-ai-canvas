// Grsai shares the WAN model family with DashScope, but not its request or task protocol.
export function grsaiWanVideoSpec({ ref, omit, coalesce, addNumbers, map, filter, eq, ne, gt, gte, lte, and, len, lower, trim, toInt, split, conditional, first, sorted, bearer, jsonCreate }) {
  const id = "grsai-wan-video", model = "wan-3.0";
  const ratios = ["16:9", "9:16", "21:9", "4:3", "1:1", "3:4"], resolutions = ["480p", "720p", "1080p"];
  const text = "text_to_video", frames = "image_to_video", reference = "reference_to_video", audio = "audio_to_video";
  const request = (path) => ref(`request.${path}`), bound = (name) => ref(`bindings.${name}`);
  const or = (...items) => ({ $or: items }), not = (value) => ({ $not: value }), oneOf = (value, items) => ({ $in: [value, items] });
  const all = (items, test) => eq(len(filter(items, "media", not(test))), 0);
  const check = (assert, message) => ({ assert, message: `Grsai Wan · ${message}` });
  const present = (value) => gt(len(trim(value)), 0);
  const values = (items) => map(items, "media", ref("media.value"));
  const role = ref("media.role"), kind = ref("media.kind"), value = ref("media.value"), mime = lower(trim(ref("media.mimeType")));
  const hasRole = (items, roles) => filter(items, "media", oneOf(role, roles));
  const imageRoles = ["", "reference", "reference_image", "subject_reference", "style_reference"];
  const images = bound("images"), videos = bound("videos"), audios = bound("audios");
  const firstFrames = bound("firstFrames"), lastFrames = bound("lastFrames"), refs = bound("referenceImages");
  const hasFrames = gt(addNumbers(len(firstFrames), len(lastFrames)), 0), op = bound("operation");
  const typed = { $concatArrays: [request("images"), request("videos"), request("audios")] };
  // IDs identify hydrated assets; anonymous inputs fall back to their actual media source.
  // The ordered fast path avoids repeated scans of large inline media for ordinary host requests.
  const identities = (items) => map(items, "media", { kind, role, id: ref("media.id"), source: conditional(not(present(ref("media.id"))), value) });
  const inputIDs = bound("inputIDs"), typedIDs = bound("typedIDs");
  const copies = (ids) => len(filter(ids, "other", eq(ref("candidate"), ref("other"))));
  const sameInputs = conditional(eq(len(typedIDs), 0), true, conditional(eq(inputIDs, typedIDs), true,
    and(eq(len(inputIDs), len(typedIDs)), eq(len(filter(inputIDs, "candidate", ne(copies(inputIDs), copies(typedIDs)))), 0))));
  const byKind = (type) => sorted(conditional(gt(len(typed), 0), request(`${type}s`), filter(request("inputs"), "media", eq(kind, type))));
  const bindings = { inputIDs: identities(request("inputs")), typedIDs: identities(typed), images: byKind("image"), videos: byKind("video"), audios: byKind("audio"),
    firstFrames: hasRole(images, ["first_frame"]), lastFrames: hasRole(images, ["last_frame"]), referenceImages: hasRole(images, imageRoles),
    operation: coalesce(request("operation"), conditional(hasFrames, frames, conditional(gt(addNumbers(len(images), len(videos)), 0), reference, conditional(gt(len(audios), 0), audio, text))))
  };
  const output = (field) => request(`output.${field}`), durationMs = ref("media.metadata.durationMs");
  const options = request(`providerOptions.${id}`), seed = request(`providerOptions.${id}.seed`);
  // Match only the small media type header, never a regular expression over the Base64 payload.
  const dataParts = split(value, ",");
  const dataURI = (type) => conditional(eq(len(dataParts), 2),
    and(eq({ $countMatches: [first(dataParts), `^data:${type};base64$`] }, 1), gt(len({ $dataPayload: value }), 0)), false);
  const publicURL = (url) => ({ $publicHttpsUrl: url });
  const mp4URL = { $countMatches: [value, "(?i)^https://[^?#]+\\.mp4(?:[?#].*)?$"] };
  const validMP4 = and(or(eq(mime, ""), eq(mime, "video/mp4")),
    conditional(publicURL(value), or(eq(mime, "video/mp4"), mp4URL), dataURI("video/mp4")));
  const mediaSource = (type) => and(or(eq(mime, ""), eq(mime, type), eq(first(split(mime, "/")), type)),
    conditional(publicURL(value), true, dataURI(`${type}/[A-Za-z0-9.+-]+`)));
  const resultID = ref("response.id"), expectedID = ref("taskId"), status = ref("response.status"), progress = ref("response.progress");
  const hasID = and(eq(resultID, { $toString: resultID }), present(resultID));
  const validID = and(hasID, or(not(present(expectedID)), eq(expectedID, resultID)));
  const results = ref("response.results"), resultURLs = map(results, "item", ref("item.url"));
  const validResults = and(eq(results, map(results, "item", ref("item"))), gt(len(resultURLs), 0), eq(len(resultURLs), len(results)),
    eq(len(filter(resultURLs, "u", not(publicURL(ref("u"))))), 0));
  const success = and(validID, eq(status, "succeeded"), not(present(ref("response.error"))), validResults);
  return {
    id, providerId: id, version: "1.0.0", name: "Grsai Wan Video", vendor: "Grsai", author: "影策社区", capability: "video",
    readmeIntro: "影策社区维护的 Grsai WAN 3.0 专属视频请求协议，非 Grsai 官方发行。安装本包后在渠道模型的请求协议选择 Grsai Wan Video，上游模型 ID 填 wan-3.0；Base URL 使用 https://grsaiapi.com 或 https://grsai.dakka.com.cn，在渠道凭据中配置 API Key。",
    baseUrl: "https://grsaiapi.com", auth: bearer, detachedPolling: true,
    videoModes: { [model]: {
      [text]: { maxImages: 0, inputRoles: [], ratios },
      [frames]: { maxImages: 2, inputRoles: ["first_frame", "last_frame"], ratios },
      [reference]: { maxImages: 10, inputRoles: ["reference_image", "reference_video", "reference_audio"], ratios },
      [audio]: { maxImages: 0, inputRoles: ["reference_audio"], ratios }
    } },
    videoOutputs: { [model]: { resolutions, minSeconds: 2, maxSeconds: 30, generateAudio: false } },
    params: [
      ["model", "string", true, "model", "仅 wan-3.0。"], ["prompt", "string", true, "prompt", "不可为空。"],
      ["duration", "integer", true, "duration", "生成时长 2–30 秒。"], ["aspectRatio", "string", true, "aspectRatio", ratios.join("、")],
      ["resolution", "string", true, "resolution", resolutions.join("、")],
      ["images", "media[]", false, "images[] / firstFrame / lastFrame", "参考图最多 10 张；显式首尾帧独立映射，不重复出现在 images 中。"],
      ["videos", "media[]", false, "videos[].url", "最多 5 段 MP4，总时长不超过 15 秒。"],
      ["audios", "media[]", false, "audios[]", "最多 5 段，每段最长 15 秒。"],
      ["providerOptions", "object", false, "seed", "仅接收本插件命名空间的整数 seed，包括 0。"]
    ],
    requestBindings: Object.entries(bindings).map(([name, value]) => ({ name, value })),
    validations: [
      check(eq(request("model"), model), "仅支持 wan-3.0。"), check(present(request("prompt")), "提示词不能为空。"),
      check(oneOf(output("aspectRatio"), ratios), "请选择支持的画幅比例。"), check(oneOf(output("resolution"), resolutions), "分辨率须为 480p、720p 或 1080p。"),
      check(and(gte(output("duration"), 2), lte(output("duration"), 30)), "输出时长须为 2–30 秒。"),
      check(or(eq(output("count"), null), eq(output("count"), 1)), "每个任务生成一个视频。"),
      check(and(not(output("watermark")), not(output("generateAudio"))), "上游未提供水印或生成音轨开关。"),
      check(or(eq(options, null), eq(options, { seed })), "不支持 seed 以外的扩展参数或空值。"),
      check(or(eq(seed, null), eq(seed, toInt(seed))), "seed 必须为整数。"),
      check(oneOf(op, [text, frames, reference, audio]), "仅支持文生、关键帧、全能参考或音频参考生成模式。"),
      check(lte(len(request("inputs")), 20), "参考素材数量超过上限。"),
      check(sameInputs, "两份素材列表不一致。"),
      check(eq(len(request("inputs")), addNumbers(len(images), len(videos), len(audios))), "仅支持图片、视频与音频素材。"),
      check(all(images, and(oneOf(kind, ["", "image"]), oneOf(role, [...imageRoles, "first_frame", "last_frame"]))), "图片角色无效。"),
      check(all(videos, and(oneOf(kind, ["", "video"]), oneOf(role, ["", "reference", "reference_video"]))), "视频角色无效。"),
      check(all(audios, and(oneOf(kind, ["", "audio"]), oneOf(role, ["", "reference", "reference_audio", "reference_voice"]))), "音频角色无效。"),
      check(and(lte(len(firstFrames), 1), lte(len(lastFrames), 1), lte(len(refs), 10)), "首尾帧各最多 1 张，普通参考图最多 10 张。"),
      check(or(not(hasFrames), eq(addNumbers(len(refs), len(audios)), 0)), "首尾帧与参考图、参考音频互斥。"),
      check(or(not(hasFrames), and(eq(len(firstFrames), 1), eq(len(videos), 0))), "本插件暂不开放单独尾帧或首尾帧混合参考视频。"),
      check(and(lte(len(videos), 5), lte(len(audios), 5)), "参考视频和音频分别最多 5 段。"),
      check(all(images, mediaSource("image")), "图片须使用公网 HTTPS URL 或图片 Base64 data URI。"),
      check(all(audios, mediaSource("audio")), "音频须使用公网 HTTPS URL 或音频 Base64 data URI。"),
      check(all(videos, validMP4), "参考视频仅支持 MP4；请提供 video/mp4 或可识别的 MP4 地址。"),
      check(and(all(videos, gt(durationMs, 0)), lte({ $sum: map(videos, "media", durationMs) }, 15000)), "参考视频须有时长，总时长不超过 15 秒。"),
      check(all(audios, and(gt(durationMs, 0), lte(durationMs, 15000))), "参考音频须有时长，每段不超过 15 秒。"),
      check(or(ne(op, text), eq(len(request("inputs")), 0)), "文生视频不能携带素材。"),
      check(or(ne(op, frames), hasFrames), "图生视频须提供首帧。"),
      check(or(ne(op, reference), and(not(hasFrames), gt(len(request("inputs")), 0))), "全能参考须提供普通参考素材，不能使用首尾帧。"),
      check(or(ne(op, audio), and(gt(len(audios), 0), eq(addNumbers(len(images), len(videos)), 0))), "音频参考模式只接受参考音频。")
    ],
    create: jsonCreate("/v1/api/generate", {
      model: request("model"), prompt: request("prompt"), aspectRatio: output("aspectRatio"), resolution: output("resolution"), duration: output("duration"),
      replyType: "async", seed: omit(seed), firstFrame: first(values(firstFrames)), lastFrame: first(values(lastFrames)),
      images: omit(values(refs)), videos: omit(map(videos, "media", { url: value })), audios: omit(values(audios))
    }, { originPath: true }),
    poll: { method: "GET", path: "/v1/api/result", originPath: true, query: { id: expectedID } },
    response: {
      // Preserve the received ID: the host rejects mismatches before a different job can affect this order.
      taskId: conditional(hasID, resultID),
      status: conditional(success, "succeeded", conditional(and(validID, eq(status, "running")), "processing", "failed")),
      message: coalesce(ref("response.error"), { $switch: { cases: [
        { when: not(validID), then: "Grsai 返回了空或不匹配的任务 ID" },
        { when: eq(status, "violation"), then: "Grsai 内容审核未通过" },
        { when: eq(status, "failed"), then: "Grsai 视频生成失败" },
        { when: and(eq(status, "succeeded"), not(validResults)), then: "Grsai 完成响应缺少有效视频地址" },
        { when: not(oneOf(status, ["succeeded", "running"])), then: "Grsai 返回未知任务状态" }
      ] } }),
      stage: conditional(and(eq(status, "running"), ne(progress, null), eq(progress, toInt(progress)), gte(progress, 0), lte(progress, 100)), { $concat: [progress, "%"] }),
      videos: conditional(success, resultURLs), usage: ref("response.usage"), errorPaths: ["error"], resultEphemeral: true
    },
    notes: "仅适配 [Grsai Wan](https://qmy27nhsd9.apifox.cn/531060669e0) 和[异步结果查询](https://qmy27nhsd9.apifox.cn/452409577e0)，并非 DashScope、RollDek 或 OpenAI 视频协议。全球节点为 https://grsaiapi.com，国内节点为 https://grsai.dakka.com.cn；渠道任选其一，同一任务在原节点查询。固定 replyType=async，不提供流式输出、回调、上游取消或模型列表。running 回执不是成功；成功读取全部 results[].url，违规和失败均终止。结果生成后仅 2 小时内能查询，不等于文件 URL 的有效期；宿主应持久化结果。可选 usage 原样保留，查询页未承诺返回 usage，不能把 credits 当 Token。供应商费用按（生成秒数＋参考视频秒数）×输出分辨率对应秒价；文档没有单价或小数取整规则。\n\n普通参考图、音频使用字符串数组，视频使用 {url} 对象；素材按 order 排序。音频按单段最长 15 秒解释，视频按总和 15 秒；插件要求音视频带可信时长信息，拒绝无法检查的输入。宿主仅提交公网 HTTPS URL 或同类别的 Base64 data URI；不接收私网、HTTP 或无类型裸 Base64。上游明确首尾帧与 images、audios 互斥；示例混传不能作为有效组合。单独尾帧、首尾帧搭配参考视频未得到明确确认，本版主动收窄，不推断厂商禁用。未声明上游参考图编号语法。 纯音频参考模式表示“提示词＋参考音频”，不是无提示词；字段未要求必须附带画面素材，但纯音频组合需上游实测。视频编辑、增强等不在本接口范围；宿主的 video_to_video 价格选择器不是此处的请求模式。"
  };
}
