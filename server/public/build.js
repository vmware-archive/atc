var autoscroll = false;

function streamLog(uri) {
  var ws = new WebSocket(uri);

  ws.onmessage = function(event) {
    $("#build-log").append(event.data);

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
  var title = $("#build-title");

  if (title.hasClass("pending") || title.hasClass("started"))
    autoscroll = true;

  $(window).scroll(function() {
    var scrollEnd = $(window).scrollTop() + $(window).height();

    if (scrollEnd >= ($(document).height() - 16)) {
      autoscroll = true;
    } else {
      autoscroll = false;
    }
  });

  $("#builds").bind('mousewheel', function(e){
    if (e.originalEvent.deltaX != 0) {
      $(this).scrollLeft($(this).scrollLeft() + e.originalEvent.deltaX);
    } else {
      $(this).scrollLeft($(this).scrollLeft() - e.originalEvent.deltaY);
    }

    return false;
  });

  if ($(".build-metadata").size() > 1)
    $(".build-metadata").hide();

  $(".resource-header").click(function() {
    $(this).parent().find(".build-metadata").toggle();
  });

  scrollToCurrentBuild();
});
