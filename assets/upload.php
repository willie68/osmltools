
<?php
require_once 'Database.php';

try {
    $db = new Database();

    if ($_SERVER['REQUEST_METHOD'] === 'POST' && isset($_FILES['datei']) && isset($_POST['vesselid']) && isset($_POST['username'])) {
        $vesselId = trim($_POST['vesselid']);
        $username = trim($_POST['username']);

        if (!ctype_digit($vesselId) || (int)$vesselId < 0 || (int)$vesselId > 65535) {
            throw new Exception("Ungültige Fahrzeug-ID. Bitte eine ganze Zahl zwischen 0 und 65535 eingeben.");
        }

        $uploadDir = __DIR__ . '/uploads/';
        if (!is_dir($uploadDir) && !mkdir($uploadDir, 0777, true)) {
            throw new Exception("Upload-Verzeichnis konnte nicht erstellt werden.");
        }

        $fileTmpPath = $_FILES['datei']['tmp_name'];
        $fileName = basename($_FILES['datei']['name']);
        $destination = $uploadDir . $fileName;

        if (!move_uploaded_file($fileTmpPath, $destination)) {
            throw new Exception("Fehler beim Verschieben der hochgeladenen Datei.");
        }

        $fileDate = date('Y-m-d H:i:s');
        $sha256 = hash_file('sha256', $destination);

        $success = $db->insertUpload($fileName, $fileDate, $sha256, (int)$vesselId, $username);
        if ($success) {
            echo "✅ Datei erfolgreich hochgeladen und gespeichert.";
        } else {
            echo "❌ Fehler beim Speichern in der Datenbank.";
        }
    }
} catch (Exception $e) {
    echo "❌ Fehler: " . htmlspecialchars($e->getMessage());
}
?>

<form method="post" enctype="multipart/form-data">
    <label for="datei">Datei auswählen:</label>
    <input type="file" name="datei" required><br><br>

    <label for="vesselid">Fahrzeug-ID (0–65535):</label>
    <input type="number" name="vesselid" min="0" max="65535" required><br><br>

    <label for="username">Benutzer:</label>
    <input type="text" name="username" required><br><br>

    <button type="submit">Hochladen</button>
</form>
