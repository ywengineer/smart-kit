package loaders

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestViperYamlLoader(t *testing.T) {
	//
	loader := NewViperLoader("conf_viper", Yaml)
	var c = &Conf{}
	if err := loader.Load(c); err != nil {
		t.Fatalf("%v", err)
	} else {
		t.Log(sonic.MarshalString(c))
	}
}

func TestViper(t *testing.T) {

	// 1. 初始化Viper实例
	v := viper.New()

	// 2. 配置YAML文件读取参数
	v.SetConfigName("conf_viper") // 配置文件名称(不含扩展名)
	v.SetConfigType("yaml")       // 配置文件类型
	v.AddConfigPath(".")          // 配置文件所在路径(当前目录)
	v.AddConfigPath("./configs")  // 可选：增加其他可能的配置文件路径

	// 3. 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("❌ 读取配置文件失败: %v\n", err)
		// 判断是否是配置文件未找到错误
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			fmt.Println("   请确保配置文件存在于指定路径")
		}
		os.Exit(1)
	}
	fmt.Println("✅ 配置文件读取成功")

	// 4. 配置环境变量解析（关键步骤，支持Windows环境变量）
	v.AutomaticEnv() // 自动绑定环境变量
	// 设置环境变量与配置键的映射规则（处理横杠和驼峰命名）
	// 例如：配置中的load-balance会映射到环境变量LOAD_BALANCE
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	// 4. 调试：打印原始值，确认占位符是否被解析
	workersRaw := v.Get("workers")
	fmt.Printf("调试：workers 原始值 = %v（类型：%T）\n", workersRaw, workersRaw)

	// 5. 手动处理 workers 解析（容错处理）
	var workers int
	switch val := workersRaw.(type) {
	case string:
		// 若仍为字符串（未解析），尝试提取默认值
		if strings.HasPrefix(val, "${") && strings.Contains(val, ":") {
			// 从 ${APP_WORKERS:10} 中提取默认值 10
			defaultVal := strings.Split(strings.Trim(val, "{}"), ":")[1]
			parsed, err := strconv.Atoi(defaultVal)
			if err != nil {
				fmt.Printf("解析默认值失败: %v\n", err)
				os.Exit(1)
			}
			workers = parsed
			fmt.Printf("使用默认值: workers = %d\n", workers)
		} else {
			fmt.Printf("workers 值无效: %s\n", val)
			os.Exit(1)
		}
	case int:
		workers = val
	default:
		fmt.Printf("workers 类型不支持: %T\n", val)
		os.Exit(1)
	}

	// 5. 解析配置到结构体
	var cfg Conf
	if err := v.Unmarshal(&cfg); err != nil {
		fmt.Printf("❌ 配置解析失败: %v\n", err)
		os.Exit(1)
	}

	// 6. 打印解析结果
	fmt.Println("\n📊 解析后的配置:")
	fmt.Printf("  网络类型: %s\n", cfg.Network)
	fmt.Printf("  监听地址: %s\n", cfg.Address)
	fmt.Printf("  工作线程数: %d\n", cfg.Workers)
	fmt.Printf("  负载均衡策略: %s\n", cfg.WorkerLoadBalance)
	fmt.Printf("  服务名称: %s\n", cfg.ServiceName)
	fmt.Printf("  服务权重: %d\n", cfg.Weight)
	fmt.Println("  元数据:")
	for k, v := range cfg.Metadata {
		fmt.Printf("    %s: %s\n", k, v)
	}
}
