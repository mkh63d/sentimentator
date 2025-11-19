import Swal from 'sweetalert2';

export function fireTimedAlert(title: string, message: string, timer: number, showLoading = true) {
  let timerInterval: ReturnType<typeof setInterval> | undefined;
  Swal.fire({
    title: title,
    html: message,
    timer: timer,
    timerProgressBar: true,
    didOpen: () => {
      if (showLoading) {
        Swal.showLoading();
      }
      const timerEl = Swal.getPopup()?.querySelector('b');
      timerInterval = setInterval(() => {
        if (timerEl) {
          timerEl.textContent = `${Swal.getTimerLeft()}`;
        }
      }, 100);
    },
    willClose: () => {
      if (timerInterval !== undefined) {
        clearInterval(timerInterval);
      }
    },
  }).then((result) => {
    if (result.dismiss === Swal.DismissReason.timer) {
      console.log('Alert closed by timer');
    }
  });
}

// export function fireInputAlert(title: string, message: string, confirmButtonText = 'Submit') {
//   Swal.fire({
//     title: title,
//     input: 'text',
//     inputAttributes: {
//       autocapitalize: 'off',
//     },
//     showCancelButton: true,
//     confirmButtonText: confirmButtonText,
//     showLoaderOnConfirm: true,
//     preConfirm: async (inputValue) => {
//       try {

//       }
//   })
// }
