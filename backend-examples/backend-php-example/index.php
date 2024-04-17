<?php
require 'vendor/autoload.php';

use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Slim\Factory\AppFactory;
use GuzzleHttp\Client;

$app = AppFactory::create();

// Add CORS Middleware
$app->add(function ($request, $handler) {
    $response = $handler->handle($request);
    return $response
            ->withHeader('Access-Control-Allow-Origin', '*')
            ->withHeader('Access-Control-Allow-Headers', '*')
            ->withHeader('Access-Control-Allow-Methods', '*');
});

// Read PEM file function
function readPemFile($filePath) {
    return file_get_contents($filePath);
}

// Sign message function
function sign($message, $privateKeyStr) {
    // Load the private key
    $privateKey = openssl_pkey_get_private($privateKeyStr);
    if (!$privateKey) {
        throw new Exception('Private key is not valid');
    }

    // Generate the signature
    $signature = '';
    openssl_sign($message, $signature, $privateKey, OPENSSL_ALGO_SHA256);

    // Encode the signature in base64 to make it readable
    $base64Signature = base64_encode($signature);

    // Free the private key from memory
    openssl_free_key($privateKey);

    return $base64Signature;
}

// Post authentication details function
function postAuthenticationDetails($companyCode, $userId, $signature, $timestamp) {
    $client = new Client([
        'verify' => false
    ]);
    $url = 'https://stylist-auth-api-b44moh36lq-ey.a.run.app/api/authenticate';
    $body = [
        'company_code' => $companyCode,
        'user_id' => $userId,
        'signature' => $signature,
        'timestamp' => $timestamp,
    ];
    $response = $client->post($url, ['json' => $body, 'headers' => ['Content-Type' => 'application/json']]);
    $data = json_decode($response->getBody(), true);
    return $data['token'];
}

// Authentication route
$app->post('/demo/v1/authenticate', function (Request $request, Response $response) {
    // mb_internal_encoding('UTF-8');
    $data = json_decode($request->getBody(), true);
    $userId = $data['user_id'] ?? null;
    $companyCode = $data['company_code'] ?? null;

    if (!$userId || !$companyCode) {
        $response->getBody()->write(json_encode(['error' => 'Missing user_id or company_code']));
        return $response->withStatus(400)->withHeader('Content-Type', 'application/json');
    }

    $privateKeyStr = readPemFile("$companyCode.cer");
    $timestamp = date('Y-m-d H:i:s');
    $signedMessage = sign($timestamp, $privateKeyStr);

    $token = postAuthenticationDetails($companyCode, $userId, $signedMessage, $timestamp);

    $response->getBody()->write(json_encode(['token' => $token]));
    return $response->withHeader('Content-Type', 'application/json');
});

$app->run();
