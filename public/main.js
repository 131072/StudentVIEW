navigator.serviceWorker && navigator.serviceWorker.register('/studentview/sw.js').then(function(registration) {
  console.log('Excellent, registered with scope: ', registration.scope);
});

function newgraph(x, y, z) {
  var bar = new ProgressBar.SemiCircle(x, {
    strokeWidth: 3,
    color: '#FFEA82',
    trailColor: '#eee',
    trailWidth: 1,
    easing: 'easeInOut',
    duration: 2800,
    svgStyle: null,
    text: {
      value: '',
      alignToBottom: false
    },
    from: {color: 'rgb(150, 0, 0)'},
    to: {color: 'rgb(0, 168, 0)'},
    // Set default step function for all animate calls
    step: (state, bar) => {
      bar.path.setAttribute('stroke', state.color);
      var temp = bar.value() * 100
      var value = parseFloat(temp.toFixed(1))
      if (value === 0) {
        bar.setText('');
      } else {
        if (z == undefined) {
          bar.setText(value + "%")
        }
        else {
          bar.setText(value + "% (" + svue.toLetter(value) + ")");
        }
      }

      bar.text.style.color = state.color;
    }
  });
  bar.text.style.fontFamily = '"Raleway", Helvetica, sans-serif';
  bar.text.style.fontSize = '100%';
  x.className = y;
  bar.setText(z)
  return bar
}

function scrollTo(x) {
  console.log(x)
  console.log($(window).scrollTo(x));
}
