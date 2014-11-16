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

var v1Handlers = {
  "log": function(msg) {
    processLogs(JSON.parse(msg.data));
  },

  "error": function(msg) {
    if (msg.data === undefined) {
      // 'error' event may also be native browser error, unfortunately
      return
    }

    processError(JSON.parse(msg.data));
  },

  "status": function(msg) {
    var event = JSON.parse(msg.data);

    var currentStatus = $("#page-header").attr("class");

    var buildTimes = $(".build-times");

    var status = event.status;
    var m = moment.unix(event.time);

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

      var startTime = $(".build-times time.started").attr("datetime");

      // Some events cause the build to never start (e.g. input errors).
      var didStart = !!startTime

      if(didStart) {
        var duration = moment.duration(m.diff(moment(startTime)));

        var durationEle = $("<span>");
        durationEle.addClass("duration");
        durationEle.text(duration.format("h[h]m[m]s[s]"));

        $("<dt/>").text("duration").appendTo(buildTimes);
        $("<dd/>").append(durationEle).appendTo(buildTimes);
      }
    }

    // only transition from transient states; state may already be set
    // if the page loaded after build was done
    if(currentStatus != "pending" && currentStatus != "started") {
      return;
    }

    $("#page-header").attr("class", status);
    $("#builds .current").attr("class", status + " current");

    if(status != "started") {
      $(".abort-build").remove();
    }
  },

  "input": function(msg) {
    renderResource(JSON.parse(msg.data), msg.type);
  },

  "output": function(msg) {
    renderResource(JSON.parse(msg.data), msg.type);
  }
}

var eventHandlers = {
  "0.0": {
    "log": function(msg) {
      writeLogs(msg.data, $("#build-log"));
    },
  },

  "1.0": v1Handlers,

  "1.1": v1Handlers,
}

function renderResource(event, type) {
  var resource = event[type];
  var info = resourceInfo(resource.resource, type);

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

  if (autoscroll) {
    $(document).scrollTop($(document).height());
  }
}

function streamLog(uri) {
  var es = new EventSource(uri);

  var successfullyConnected = false;
  var eventHandler;
  var currentVersion;

  es.addEventListener("version", function(event) {
    successfullyConnected = true;

    if (eventHandler) {
      for (var key in eventHandler) {
        es.removeEventListener(key, eventHandler[key]);
      }
    }

    currentVersion = JSON.parse(event.data);
    eventHandler = eventHandlers[currentVersion];

    for (var key in eventHandler) {
      es.addEventListener(key, eventHandler[key], false);
    }
  });

  es.addEventListener("end", function(event) {
    es.close();
  });

  es.onerror = function(event) {
    if(currentVersion != "1.1") {
      // versions < 1.1 cannot distinguish between end of stream and an
      // interrupted connection
      es.close();
    }

    if(!successfullyConnected) {
      // assume rejected because of auth
      $("#build-requires-auth").show();
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
  var title = $("#page-header");

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

  $(".build-source .header").click(function() {
    $(this).parent().find(".resource-body").toggle();
  });

  scrollToCurrentBuild();
});
