function handleFormError() {
  document.addEventListener("htmx:responseError", (event) => {
    const xhr = event.detail.xhr;
    const errors = JSON.parse(xhr.responseText);
    console.log("errors", errors);

    if (xhr.status === 422) {
      for (const formId of Object.keys(errors)) {
        const formErrors = errors[formId];

        for (const name of Object.keys(formErrors)) {
          const field = document.querySelector(
            `#${formId} [data-for="${name}"]`,
          );

          console.log("field", field);
          if (!field) return;

          field.setAttribute("data-error", formErrors[name]);
        }
      }
    } else if (xhr.status === 404) {
      alert("Data is not found");
    } else if (xhr.status === 400) {
      alert(errors.error);
    }
  });
}
