
<?php
require_once 'Database.php';

try {
    $db = new Database();
    $vesselIds = $db->getDistinctVesselIds();

    $filterId = isset($_GET['filter_vesselid']) ? $_GET['filter_vesselid'] : null;
    if ($filterId !== null && ctype_digit($filterId) && (int)$filterId >= 0 && (int)$filterId <= 65535) {
        $uploads = $db->getUploadsByVesselId((int)$filterId);
    } else {
        $uploads = $db->getAllUploads();
    }
} catch (Exception $e) {
    echo "❌ Fehler beim Laden der Daten: " . htmlspecialchars($e->getMessage());
    $uploads = [];
    $vesselIds = [];
}
?>

<form method="get">
    <label for="filter_vesselid">Fahrzeug-ID filtern:</label>
    <select name="filter_vesselid" onchange="this.form.submit()">
        <option value="">-- Alle --</option>
        <?php foreach ($vesselIds as $row): ?>
            <option value="<?= htmlspecialchars($row['vesselid']) ?>"
                <?= ($filterId == $row['vesselid']) ? 'selected' : '' ?>>
                <?= htmlspecialchars($row['vesselid']) ?>
            </option>
        <?php endforeach; ?>
    </select>
</form>

<table border="1" cellpadding="5" cellspacing="0">
    <thead>
        <tr>
            <th>ID</th>
            <th>Dateiname</th>
            <th>Upload-Datum</th>
            <th>SHA256 Hash</th>
            <th>Fahrzeug-ID</th>
            <th>Benutzer</th>
        </tr>
    </thead>
    <tbody>
        <?php if (empty($uploads)): ?>
            <tr><td colspan="6">Keine Einträge gefunden.</td></tr>
        <?php else: ?>
            <?php foreach ($uploads as $upload): ?>
                <tr>
                    <td><?= htmlspecialchars($upload['id']) ?></td>
                    <td><?= htmlspecialchars($upload['filename']) ?></td>
                    <td><?= htmlspecialchars($upload['filedate']) ?></td>
                    <td><?= htmlspecialchars($upload['sha256hash']) ?></td>
                    <td><?= htmlspecialchars($upload['vesselid']) ?></td>
                    <td><?= htmlspecialchars($upload['username']) ?></td>
                </tr>
            <?php endforeach; ?>
        <?php endif; ?>
    </tbody>
</table>
