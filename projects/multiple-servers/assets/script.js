function fetchImages() {
  return fetch("http://localhost:8081/api/images.json").then((_) => _.json());
}

function timeout(t, v) {
  return new Promise((res) => {
    setTimeout(() => res(v), t);
  });
}

const gallery$ = document.querySelector(".gallery");

fetchImages().then(
  (images) => {
    console.log(images);
    gallery$.textContent = images.length ? "" : "No images available.";

    images.forEach((img) => {
      const imgElem$ = document.createElement("img");
      imgElem$.src = img.URL;
      imgElem$.alt = img.Alt_Text;
      const titleElem$ = document.createElement("h3");
      titleElem$.textContent = img.Title;
      const wrapperElem$ = document.createElement("div");
      wrapperElem$.classList.add("gallery-image");
      wrapperElem$.appendChild(titleElem$);
      wrapperElem$.appendChild(imgElem$);
      gallery$.appendChild(wrapperElem$);
    });
  },
  () => {
    gallery$.textContent = "Something went wrong.";
  }
);
