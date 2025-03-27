// 引入 crypto-js 库用于 HMAC - SHA256 加密
const CryptoJS = require('crypto-js');
// 生成签名的函数
function generateSignature(params, apiKey) {
    // 拼接参数名和参数值
    let paramStr = Object.keys(params)
        .sort()
        .map(key => {
            let ek = encodeURIComponent(key)
            let ev = encodeURIComponent(params[key])
            return `${ek}=${ev}`
        })
        .join("&")
    // 使用 HMAC - SHA256 进行加密
    const hmac = CryptoJS.HmacSHA256(paramStr, apiKey);
    // 将加密结果转换为十六进制字符串
    return hmac.toString(CryptoJS.enc.Hex);
}