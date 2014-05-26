var autoscroll = true;

function streamLog(uri) {
  var ws = new WebSocket(uri);

  ws.onmessage = function(event) {
    document.getElementById("build-log").innerHTML += event.data

    if (autoscroll) {
      $(document).scrollTop($(document).height());
    }
  };
}

$(document).ready(function() {
  $(window).scroll(function() {
    var scrollEnd = $(window).scrollTop() + $(window).height();

    if (scrollEnd >= ($(document).height() - 16)) {
      autoscroll = true;
    } else {
      autoscroll = false;
    }

    scrolled = true;
  });
});
