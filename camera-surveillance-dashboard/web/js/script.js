
function triggerScrollingEffect() {
  const topNavigationBar = document.querySelector(".topNavigationBar");

  window.addEventListener("scroll", () => {
    if (window.scrollY > 90) {
      topNavigationBar.classList.add("bg-dark");
      topNavigationBar.classList.add("add-opactity");
    } else {
      topNavigationBar.classList.remove("bg-dark");
      topNavigationBar.classList.remove("add-opactity");
    }
  });
}

function increaseReadings() {
  let readings = document.querySelectorAll(".reading");
  readings.forEach((reading) => {
    reading.innerText = 0;

    let updateReading = () => {
      let data_target = +reading.getAttribute("data-target");
      // console.log("individualReading is:" + data_target);
      let individualReading = +reading.innerText;

      let incrementValue = data_target / 100;
      if (individualReading < data_target) {
        reading.innerText = Math.ceil( individualReading + incrementValue)
        //  call this function recursively 
        setTimeout(updateReading, 10)
      } else {
        reading.innerText = data_target
      }
    };

    updateReading()
  });
}

// call the triggerScrollingEffect once DOM has loaded
document.addEventListener("DOMContentLoaded", triggerScrollingEffect);

// call the increaseReadings once DOM has loaded
document.addEventListener("DOMContentLoaded", increaseReadings);
