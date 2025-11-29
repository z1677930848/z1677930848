<?php
/**
 * 统一响应处理类
 */

class Response {
    /**
     * 成功响应
     */
    public static function success($data = null, $message = 'success') {
        self::output([
            'code' => 200,
            'message' => $message,
            'data' => $data
        ]);
    }

    /**
     * 错误响应
     */
    public static function error($message, $code = 500, $data = null) {
        self::output([
            'code' => $code,
            'message' => $message,
            'data' => $data
        ]);
    }

    /**
     * 输出JSON响应
     */
    private static function output($data) {
        header('Content-Type: application/json; charset=utf-8');
        echo json_encode($data, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT);
        exit;
    }
}
