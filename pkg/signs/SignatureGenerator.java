import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.io.UnsupportedEncodingException;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.util.*;

public class SignatureGenerator {

    public static String generate(Map<String, String> params, String apiKey)
            throws NoSuchAlgorithmException, InvalidKeyException, UnsupportedEncodingException {
        // 对参数名进行排序
        final String paramStr = params.keySet()
                .stream()
                .sorted()
                .map(key -> {
                    try {
                        String encodedKey = URLEncoder.encode(key, StandardCharsets.UTF_8.name());
                        String encodedValue = URLEncoder.encode(params.getOrDefault(key, ""), StandardCharsets.UTF_8.name());
                        return encodedKey + "=" + encodedValue;
                    } catch (UnsupportedEncodingException e) {
                        // ignore
                    }
                    return "";
                })
                .filter(v -> !v.isEmpty())
                .collect(Collectors.joining("&"));
        // 创建 HMAC - SHA256 实例
        final Mac hmacSha256 = Mac.getInstance("HmacSHA256");
        hmacSha256.init(new SecretKeySpec(apiKey.getBytes(StandardCharsets.UTF_8), "HmacSHA256"));
        // 将字节数组转换为十六进制字符串
        final StringBuilder hexString = new StringBuilder();
        for (byte b : hmacSha256.doFinal(paramStr.getBytes(StandardCharsets.UTF_8))) {
            String hex = Integer.toHexString(0xff & b);
            if (hex.length() == 1) {
                hexString.append('0');
            }
            hexString.append(hex);
        }
        return hexString.toString();
    }
}
