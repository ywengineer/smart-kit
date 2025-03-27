public class SignatureGenerator
    {
        public static string Generate(Dictionary<string, string> paramsDict, string apiKey)
        {
            // 对参数名进行排序
            var sortedKeys = paramsDict.Keys.OrderBy(k => k).ToList();

            // 拼接参数名和参数值，同时进行编码
            StringBuilder paramStr = new StringBuilder();
            foreach (var key in sortedKeys)
            {
                var value = paramsDict[key];
                if (!string.IsNullOrEmpty(value))
                {
                    var encodedKey = Uri.EscapeDataString(key);
                    var encodedValue = Uri.EscapeDataString(value);
                    if (paramStr.Length > 0)
                    {
                        paramStr.Append("&");
                    }
                    paramStr.Append(encodedKey).Append("=").Append(encodedValue);
                }
            }

            // 创建 HMAC - SHA256 实例
            using (HMACSHA256 hmac = new HMACSHA256(Encoding.UTF8.GetBytes(apiKey)))
            {
                // 计算签名
                byte[] hashBytes = hmac.ComputeHash(Encoding.UTF8.GetBytes(paramStr.ToString()));

                // 将字节数组转换为十六进制字符串
                StringBuilder hexString = new StringBuilder();
                foreach (byte b in hashBytes)
                {
                    hexString.Append(b.ToString("x2"));
                }

                return hexString.ToString();
            }
        }
    }