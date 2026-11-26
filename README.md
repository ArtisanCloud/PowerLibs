# PowerLibs

作为ArtisanCloud开发Golang项目的通用库，我们把很多基础模块常用到的方法，放在了PowerLib中。

## YAML 依赖迁移

`object.OpenYMLFile` 和 `object.SaveYMLFile` 使用 `go.yaml.in/yaml/v3`。如果传入的类型使用 `yaml.Node`，请将其导入路径从 `gopkg.in/yaml.v3` 改为 `go.yaml.in/yaml/v3`。保留旧路径时，新解码器可能跳过自定义的 `UnmarshalYAML(*yaml.Node) error` 方法，直接按字段解码。

特别感谢 <a href="https://github.com/northseadl">Northseadl</a>的贡献

<a href="https://github.com/northseadl">
  <img width="50" src="https://avatars.githubusercontent.com/u/99570485?v=4">
</a>
