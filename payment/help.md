
## Protobuf tag 的格式
Protobuf 生成的 Go 结构体标签格式通常为：

protobuf:"<wire_type>,<tag>,[opt|req|rep],[name=<field_name>]"
```text
wire_type：编码类型（如 varint、bytes 等，对应 Protobuf 类型）
tag：字段标识（必须与 Protobuf 定义中的 tag 一致）
opt|req|rep：可选 / 必需 / 重复（对应 Protobuf 的 optional/required/repeated）
```
