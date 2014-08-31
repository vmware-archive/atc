var autoscroll = false;

var eventHandlers = {
  "0.0": function(msg) {
    processLogs(msg.data);
  },
  "1.0": function(msg) {
    var eventMsg = JSON.parse(msg.data);

    switch(eventMsg.type) {
    case "log":
      processLogs(eventMsg.event.payload);
      break;
    case "status":
      var status = eventMsg.event.status;

      if(status != "started") {
        $(".abort-build").remove();
      }

      $("#build-title").attr("class", status);
      $("#builds .current").attr("class", status + " current");

      break;
    }
  }
}

function processLogs(payload) {
  var sequence = ansiparse(payload);

  var log = $("#build-log");
  var ele;

  for(var i = 0; i < sequence.length; i++) {
    ele = $("<span>");
    ele.text(sequence[i].text);

    if(sequence[i].foreground) {
      ele.addClass("ansi-"+sequence[i].foreground+"-fg");
    }

    if(sequence[i].background) {
      ele.addClass("ansi-"+sequence[i].background+"-bg");
    }

    if(sequence[i].bold) {
      ele.addClass("ansi-bold");
    }

    log.append(ele);
  }

  if (autoscroll) {
    $(document).scrollTop($(document).height());
  }
}

function streamLog(uri) {
  var ws = new WebSocket(uri);

  var eventHandler;
  ws.onmessage = function(event) {
    if(!eventHandler) {
      var versionMsg = JSON.parse(event.data);

      if(versionMsg.version) {
        eventHandler = eventHandlers[versionMsg.version];
      }
    } else {
      eventHandler(event);
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
