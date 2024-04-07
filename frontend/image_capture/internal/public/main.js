const startBtn = document.getElementById("startCam");
const stopBtn = document.getElementById("stopCam");
const video = document.getElementById("video");
const captureButton = document.getElementById("capture-button");
const downloadLink = document.getElementById("download-link");
const roomID = "12";

stopBtn.disabled = true;
startBtn.addEventListener("click", startCam);
async function startCam() {
  try {
    const stream = await navigator.mediaDevices.getUserMedia({
      video: { width: 640, height: 360 },
    });
    video.srcObject = stream;
    startBtn.disabled = true;
    stopBtn.disabled = false;
  } catch (err) {
    console.error("Error starting web cam:", err);
  }
}

stopBtn.addEventListener("click", stopCam);
function stopCam() {
  stop(video.srcObject);
  startBtn.disabled = false;
  stopBtn.disabled = true;
}

captureButton.addEventListener("click", captureImage);
async function captureImage() {
  try {
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");
    const reader = new FileReader();
    const roomID = document.getElementById("roomid").innerText;

    canvas.width = video.width;
    canvas.height = video.height;

    ctx.drawImage(video, 0, 0);
    canvas.toBlob(async (blob) => {
      const formData = new FormData();
      formData.append("data", blob);

      // const response = await fetch(`http://localhost:8080/img?roomID=${roomID}`,{
      const response = await fetch(`https://sensaide-onsvm4sjba-uw.a.run.app/img?roomID=${roomID}`,{
        method: "POST",
        body: formData,
      });

      if (!response.ok){
        throw new Error(`HTTP Error! status: ${response.status}}`)
      };

      console.log(await response.text());
    }, "image/jpeg");
  } catch (error) {
    console.error("Error capturing image:", error);
  }
}

function stop(stream) {
  stream.getTracks().forEach((trk) => trk.stop());
}

function myDateFmt(date) {
  return (
    date.getFullYear() +
    String(date.getMonth() + 1).padStart(2, "0") +
    String(date.getDate()).padStart(2, "0") +
    "-" +
    String(date.getHours()).padStart(2, "0") +
    String(date.getMinutes()).padStart(2, "0") +
    String(date.getSeconds()).padStart(2, "0") +
    String(date.getMilliseconds()).padStart(3, "0") // Extract milliseconds and pad with zeros
  );
}
