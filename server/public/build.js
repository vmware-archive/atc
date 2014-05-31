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

function scrollToCurrentBuild() {
  var currentBuild = $("#builds .current");
  var buildWidth = currentBuild.width();
  var left = currentBuild.offset().left;

  if((left + buildWidth) > window.innerWidth) {
    $("#builds").scrollLeft(left - buildWidth);
  }
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

  $("#builds").bind('mousewheel', function(e){
    if(e.originalEvent.deltaX != 0) {
      $(this).scrollLeft($(this).scrollLeft() + e.originalEvent.deltaX);
    } else {
      $(this).scrollLeft($(this).scrollLeft() - e.originalEvent.deltaY);
    }

    return false;
  });

  scrollToCurrentBuild();
});
