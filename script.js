// Get all light and aircon buttons
const lightButtons = document.querySelectorAll('.room button:first-child');
const airconButtons = document.querySelectorAll('.room button:last-child');

// Function to update button state based on value
function updateButton(button, value) {
  button.classList.toggle('on', value === 1);
  button.textContent = value === 1 ? 'On' : 'Off';
}

// Function to fetch and update room light state
function updateRoomLight(roomId, url) {
  fetch(url)
    .then(response => response.json())
    .then(data => updateButton(document.getElementById(`room${roomId}-light`), data.value));
}

// Function to fetch and update room aircon light state (similar to updateRoomLight)
function updateRoomAirconLight(roomId, url) {
  fetch(url)
    .then(response => response.json())
    .then(data => updateButton(document.getElementById(`room${roomId}-aircon`), data.value));
}

// Update light status every second
setInterval(() => {
  updateRoomLight(12, 'https://spgroup24.alwaysdata.net/lights/12');
  updateRoomLight(34, 'https://spgroup24.alwaysdata.net/lights/34');
}, 1000);

// Update aircon light status every second (similar logic as updateRoomLight)
setInterval(() => {
  updateRoomAirconLight(12, 'https://spgroup24.alwaysdata.net/aircon/12');
  updateRoomAirconLight(34, 'https://spgroup24.alwaysdata.net/aircon/34');
}, 1000);
