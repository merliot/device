let target = document.getElementById('download-target')
let gpios = deployTargetGpios[target.value];

// Function to refresh all select options based on current selections
function refreshGpioOptions() {
    const selects = document.querySelectorAll('.gpio-selector');
    const usedGpios = [];

    // Find all currently selected GPIOs
    selects.forEach(select => {
        const value = select.value;
        if (value) {
            usedGpios.push(value);
        }
    });

    // Refresh options in all select dropdowns
    selects.forEach(select => {
        const currentValue = select.value;
        let optionsHtml = '<option value="">Select GPIO</option>';

        gpios.forEach(gpio => {
            // If gpio is not used or it's the current value of this select, add it as an option
            if (!usedGpios.includes(gpio) || gpio === currentValue) {
                optionsHtml += `<option value="${gpio}">${gpio}</option>`;
            }
        });

        select.innerHTML = optionsHtml;
        select.value = currentValue;  // Restore the previously selected value
    });
}

function clearGpioOptions() {
    const selects = document.querySelectorAll('.gpio-selector');
    selects.forEach(select => {
        select.value = ""
    });
}

// Attach an event listener to the download-target dropdown to adjust available GPIOs
target.addEventListener('change', function() {
    const selectedTarget = this.value;
    if (deployTargetGpios[selectedTarget]) {
        gpios = deployTargetGpios[selectedTarget];
	clearGpioOptions();
        refreshGpioOptions();
    }
});

// Attach event listeners to the gpio dropdowns to ensure no duplicate selections
document.querySelectorAll('.gpio-selector').forEach(select => {
    select.addEventListener('change', refreshGpioOptions);
});

// Initial population of the GPIO options
refreshGpioOptions();
