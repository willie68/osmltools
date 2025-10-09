
<?php
class Database {
    private $host = '<dein host mit port>';
    private $db   = '<deine sql db>';
    private $user = '<dein sql user';
    private $pass = '<dein sql passwort>';
    private $charset = 'utf8mb4';
    private $mysqli;

    public function __construct() {
        $this->mysqli = new mysqli($this->host, $this->user, $this->pass, $this->db);
        if ($this->mysqli->connect_error) {
            error_log("Datenbankverbindung fehlgeschlagen: " . $this->mysqli->connect_error);
            die("Interner Fehler bei der Datenbankverbindung.");
        }
    }

    public function insertUpload($filename, $filedate, $sha256hash, $vesselid, $username) {
        if (!is_numeric($vesselid) || $vesselid < 0 || $vesselid > 65535 || empty($filename) || empty($sha256hash)) {
            error_log("Ungültige Parameter für insertUpload");
            return false;
        }

        $stmt = $this->mysqli->prepare("INSERT INTO uploads (filename, filedate, sha256hash, vesselid, username) VALUES (?, ?, ?, ?, ?)");
        if (!$stmt) {
            error_log("Prepare fehlgeschlagen: " . $this->mysqli->error);
            return false;
        }
        $stmt->bind_param("sssis", $filename, $filedate, $sha256hash, $vesselid, $username);
        $result = $stmt->execute();
        if (!$result) {
            error_log("Insert fehlgeschlagen: " . $stmt->error);
        }
        $stmt->close();
        return $result;
    }

    public function getDistinctVesselIds() {
        $result = $this->mysqli->query("SELECT DISTINCT vesselid FROM uploads ORDER BY vesselid ASC");
        if (!$result) {
            error_log("Query fehlgeschlagen: " . $this->mysqli->error);
            return [];
        }
        return $result ? $result->fetch_all(MYSQLI_ASSOC) : [];
    }

    public function getAllUploads() {
        $result = $this->mysqli->query("SELECT * FROM uploads ORDER BY filedate DESC");
            if (!$result) {
            error_log("Query fehlgeschlagen: " . $this->mysqli->error);
            return [];
        }
        return $result ? $result->fetch_all(MYSQLI_ASSOC) : [];
    }

    public function getUploadsByVesselId($vesselid) {
        if (!is_numeric($vesselid)) return [];
        $stmt = $this->mysqli->prepare("SELECT * FROM uploads WHERE vesselid = ? ORDER BY filedate DESC");
        if (!$stmt) {
            error_log("Prepare fehlgeschlagen: " . $this->mysqli->error);
            return [];
        }
        $stmt->bind_param("d", $vesselid);
        $stmt->execute();
        $result = $stmt->get_result();
        $uploads = $result->fetch_all(MYSQLI_ASSOC);
        $stmt->close();
        return $uploads;
    }

    public function __destruct() {
        if ($this->mysqli) {
            $this->mysqli->close();
        }
    }

    public function checkHash($hash) {
        if (strlen($hash) !== 64 || !ctype_xdigit($hash)) {
            return [
                "status" => "error",
                "message" => "Ungültiger Hash"
            ];
        }

        $stmt = $this->mysqli->prepare("SELECT COUNT(*) FROM uploads WHERE sha256hash = ?");
        if (!$stmt) {
            return [
                "status" => "error",
                "message" => "Datenbankfehler bei Vorbereitung der Abfrage"
            ];
        }

        $stmt->bind_param("s", $hash);
        $stmt->execute();
        $stmt->bind_result($count);
        $stmt->fetch();
        $stmt->close();

        if ($count > 0) {
            return [
                "status" => "exists",
                "message" => "Datei bereits vorhanden"
            ];
        } else {
            return [
                "status" => "notfound",
                "message" => "Datei nicht vorhanden"
            ];
        }
    }
}
?>
