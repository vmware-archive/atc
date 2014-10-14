var autoscroll = false;

function resourceInfo(name, type) {
  var group = ".build-"+type+"s"

  var resource = $(group).find(".build-source[data-resource='"+name+"']");
  if(resource.length === 0) {
    resource = $("#resource-template .build-source").clone(true).appendTo(group);
    resource.attr("data-resource", name);
    resource.find("h3").text(name);
  }

  return resource;
}

var eventHandlers = {
  "0.0": function(msg) {
    writeLogs(msg.data, $("#build-log"));
  },

  "1.0": function(msg) {
    var eventMsg = JSON.parse(msg.data);

    switch(eventMsg.type) {
    case "log":
      processLogs(eventMsg.event);
      break;

    case "error":
      processError(eventMsg.event);
      break;

    case "input":
    case "output":
      var resource = eventMsg.event[eventMsg.type];
      var info = resourceInfo(resource.name, eventMsg.type);

      info.removeClass("running");

      var version = info.find(".version");
      if(version.children().length === 0) {
        for(var key in resource.version) {
          $("<dt/>").text(key).appendTo(version);
          $("<dd/>").text(resource.version[key]).appendTo(version);
        }
      }

      var metadata = info.find(".build-metadata");
      if(metadata.children().length === 0) {
        for(var i in resource.metadata) {
          var field = resource.metadata[i];
          $("<dt/>").text(field.name).appendTo(metadata);
          $("<dd/>").text(field.value).appendTo(metadata);
        }
      }

      break;

    case "status":
      var currentStatus = $("#build-title").attr("class");

      var buildTimes = $(".build-times");

      var status = eventMsg.event.status;
      var m = moment.unix(eventMsg.event.time);

      var time = $("<time>");
      time.text(m.fromNow());
      time.attr("datetime", m.format());
      time.attr("title", m.format("lll Z"));
      time.addClass(status);

      if(status == "started") {
        $("<dt/>").text(status).appendTo(buildTimes);
        $("<dd/>").append(time).appendTo(buildTimes);
      } else {
        $("<dt/>").text(status).appendTo(buildTimes);
        $("<dd/>").append(time).appendTo(buildTimes);

        var startTime = moment($(".build-times time.started").attr("datetime"));
        var duration = moment.duration(m.diff(startTime));

        var durationEle = $("<span>");
        durationEle.addClass("duration");
        durationEle.text(duration.format("h[h]m[m]s[s]"));

        $("<dt/>").text("duration").appendTo(buildTimes);
        $("<dd/>").append(durationEle).appendTo(buildTimes);
      }

      // only transition from transient states; state may already be set
      // if the page loaded after build was done
      if(currentStatus != "pending" && currentStatus != "started") {
        break;
      }

      $("#build-title").attr("class", status);
      $("#builds .current").attr("class", status + " current");

      if(status != "started") {
        $(".abort-build").remove();
      }

      break;
    }
  }
}

function processLogs(event) {
  var log;

  switch(event.origin.type) {
  case "run":
    log = $("#build-log");
    break;
  case "input":
  case "output":
    var resource = resourceInfo(event.origin.name, event.origin.type);
    if(!resource.hasClass("running"))
      resource.addClass("running");

    log = resource.find(".log")
  }

  if(!log) {
    return;
  }

  writeLogs(event.payload, log);
}

function processError(event) {
  var log;

  var errorSpan = $("<span/>");
  errorSpan.addClass("error");
  errorSpan.text(event.message);

  if(event.origin) {
    switch(event.origin.type) {
    case "input":
    case "output":
      var resource = resourceInfo(event.origin.name, event.origin.type);
      resource.removeClass("running");
      resource.addClass("errored");

      log = resource.find(".log");
    }
  } else {
    log = $("#build-log");
  }

  if(!log) {
    return;
  }

  log.append(errorSpan);
}

var currentLine;
var lineCursor;
var lineToOverwrite;

function overwriteLine(segment) {
  lineCursor.after(segment);

  lineCursor = $(segment);

  var textLen = segment.text().length;
  lineToOverwrite.each(function() {
    var text = $(this).text();

    if(text.length >= textLen) {
      $(this).text(text.substr(textLen));
      return false;
    } else {
      $(this).remove();
      textLen -= text.length;
    }
  });
}

function writeLogs(payload, destination) {
  var sequence = ansiparse(payload);

  var ele;
  for(var i = 0; i < sequence.length; i++) {
    ele = $("<span>");
    ele.text(sequence[i].text);

    if(sequence[i].linebreak) {
      ele.addClass("linebreak");
    }

    if(sequence[i].foreground) {
      ele.addClass("ansi-"+sequence[i].foreground+"-fg");
    }

    if(sequence[i].background) {
      ele.addClass("ansi-"+sequence[i].background+"-bg");
    }

    if(sequence[i].bold) {
      ele.addClass("ansi-bold");
    }

    if(sequence[i].cr && currentLine) {
      lineCursor = currentLine;
      lineToOverwrite = currentLine.nextAll();
    } else if(sequence[i].linebreak) {
      currentLine = ele;
      lineToOverwrite = null;
    }

    if(lineToOverwrite) {
      overwriteLine(ele);
    } else {
      destination.append(ele);
    }
  }
}

var successfullyConnected = false;
var eventIndex = 0;
var currentEvent = 0;

function streamLog(uri) {
  var ws = new WebSocket(uri);

  var eventHandler;

  ws.onerror = function(event) {
    if(!successfullyConnected) {
      // assume rejected because of auth
      $("#build-requires-auth").show();
    }
  };

  ws.onopen = function(event) {
    successfullyConnected = true;

    // reset event processing so we can catch up on reconnect
    currentEvent = 0;
  }

  ws.onclose = function(event) {
    if(successfullyConnected) {
      // reconnect
      setTimeout(function() { streamLog(uri) }, 1000);
    }
  }

  ws.onmessage = function(event) {
    currentEvent++;

    // keep track of events we've already seen; skip until we catch up
    if(currentEvent <= eventIndex) {
      return;
    } else {
      eventIndex++;
    }

    if(!eventHandler) {
      var versionMsg = JSON.parse(event.data);

      if(versionMsg.version) {
        eventHandler = eventHandlers[versionMsg.version];
      }
    } else {
      eventHandler(event);
    }

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

  if (title.hasClass("pending") || title.hasClass("started")) {
    autoscroll = true;
  }

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

  $(".resource-header").click(function() {
    $(this).parent().find(".resource-body").toggle();
  });

  scrollToCurrentBuild();
});
