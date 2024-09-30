document.addEventListener('paste', async (e) => {
    e.preventDefault();
    if (e.clipboardData.files?.length > 0) {
        console.log('file');
        const data = e.clipboardData.files
        console.log(data[0])
        // its files
    } else {
        console.log('text');
        const text = e.clipboardData.getData("text/plain");
        console.log(text);
        // its text
    }
  });
  