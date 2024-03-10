const popoverTriggerList = document.querySelectorAll(
    '[data-bs-toggle="popover"]'
);

const popoverList = [...popoverTriggerList].map(
    (popoverListObject) => new bootstrap.Popover(popoverListObject)
);