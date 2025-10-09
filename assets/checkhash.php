<?php
header('Content-Type: application/json');
require_once 'Database.php';

$db = new Database();

if (isset($_GET['checkhash'])) {
    $hash = trim($_GET['checkhash']);

    if (strlen($hash) !== 64 || !ctype_xdigit($hash)) {
        http_response_code(400);
        echo json_encode([
            "status" => "error",
            "message" => "Ungültiger Hash"
        ]);
        exit;
    }

    $stmt = $db->getConnection()->prepare("SELECT COUNT(*) FROM uploads WHERE sha256hash = ?");
    if (!$stmt) {
        http_response_code(500);
        echo json_encode([
            "status" => "error",
            "message" => "Datenbankfehler"
        ]);
        exit;
    }

    $stmt->bind_param("s", $hash);
    $stmt->execute();
    $stmt->bind_result($count);
    $stmt->fetch();
    $stmt->close();

    if ($count > 0) {
        echo json_encode([
            "status" => "exists",
            "message" => "Datei bereits vorhanden"
        ]);
    } else {
        echo json_encode([
            "status" => "notfound",
            "message" => "Datei nicht vorhanden"
        ]);
    }
    exit;
} else {
    http_response_code(400);
    echo json_encode([
        "status" => "error",
        "message" => "Parameter 'checkhash' fehlt"
    ]);
}
?>