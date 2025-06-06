// Инициализация всех тултипов
document.addEventListener('DOMContentLoaded', function () {
    // Инициализация тултипов Bootstrap, если они будут использоваться
    const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });
}); 