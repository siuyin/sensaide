// Get all light and aircon buttons
const lightButtons = document.querySelectorAll('.room button:first-child');
const airconButtons = document.querySelectorAll('.room button:last-child');

// Function to toggle button state
function toggleButton(button) {
  const currentState = button.classList.contains('on');
  button.classList.toggle('on');
  button.textContent = currentState ? 'Light: Off' : 'Light: On';  // Update light text
  // Replace with your actual light control logic using the current state (on/off)
  console.log(`Light in ${button.id.split('-')[0]} & ${button.id.split('-')[1]} is now ${currentState ? 'off' : 'on'}`);
}

// Add click event listeners to light buttons
lightButtons.forEach(button => button.addEventListener('click', () => toggleButton(button)));

// Function to toggle aircon button state (similar logic as light buttons)
function toggleAirconButton(button) {
  const currentState = button.classList.contains('on');
  button.classList.toggle('on');
  button.textContent = currentState ? 'Aircon Light: Off' : 'Aircon Light: On';
  // Replace with your actual aircon light control logic using the current state (on/off)
  console.log(`Aircon light in ${button.id.split('-')[0]} & ${button.id.split('-')[1]} is now ${currentState ? 'off' : 'on'}`);
}

// Add click event listeners to aircon buttons
airconButtons.forEach(button => button.addEventListener('click', () => toggleAirconButton(button)));
