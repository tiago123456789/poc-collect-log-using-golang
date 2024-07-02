<?php

// Autoload files using the Composer autoloader.
require_once __DIR__ . '/vendor/autoload.php';

use GuzzleHttp\Client;
use Monolog\Logger;
use Monolog\Handler\AbstractProcessingHandler;

class CustomGuzzleHandler extends AbstractProcessingHandler
{
    private $client;
    private $url;

    public function __construct(Client $client, $url, $level = Logger::DEBUG, bool $bubble = true)
    {
        $this->client = $client;
        $this->url = $url;
        parent::__construct($level, $bubble);
    }

    protected function write(array $record): void
    {
        $logData = [
            "message" => $record["message"],
            "level" => strtolower($record["level_name"]),
            "service" => $record["service"],
            "timestamp" => $record["datetime"]
        ];
        $payload = json_encode($logData);
        $this->client->post($this->url, [
            'headers' => ['Content-Type' => 'application/json'],
            'body' => $payload,
        ]);
    }
}

// Create a Guzzle client
$client = new Client();

// Define the URL to which logs will be sent
$logUrl = 'http://localhost:3000/';

// Create a CustomGuzzleHandler instance
$guzzleHandler = new CustomGuzzleHandler($client, $logUrl);

// Create a logger instance
$logger = new Logger('my_logger');

// Add the handler to the logger
$logger->pushHandler($guzzleHandler);

// Create a processor to add default values
$defaultProcessor = function ($record) {
    $record['service'] = 'php-client-grpc';
    return $record;
};

// Add the processor to the logger
$logger->pushProcessor($defaultProcessor);

// Add records to the log
$logger->info('This is an info log message');
$logger->warning('This is a warning log message');
$logger->error('This is an error log message');

$logger->info('ABC DEF GJI MORE');
