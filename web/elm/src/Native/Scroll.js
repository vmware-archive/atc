var _concourse$atc$Native_Scroll = function() {
  function toBottom(id) {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      var ele = document.getElementById(id);
      ele.scrollTop = ele.scrollHeight - ele.clientHeight;
      callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
    });
  }

  function toWindowTop() {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      window.scrollTo(0, 0);
      callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
    });
  }

  function toWindowBottom() {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      window.scrollTo(0, document.body.scrollHeight);
      callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
    });
  }

  function scrollElement(id, delta) {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      document.getElementById(id).scrollLeft -= delta;
      callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
    });
  }

  function scrollIntoView(selector) {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      document.querySelector(selector).scrollIntoView();
      callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
    });
  }

  function scrollUp() {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      window.scrollBy(0, -60);
      callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
    });
  }

  function scrollDown() {
    return _elm_lang$core$Native_Scheduler.nativeBinding(function(callback) {
      window.scrollBy(0, 60);
      callback(_elm_lang$core$Native_Scheduler.succeed(_elm_lang$core$Native_Utils.Tuple0));
    });
  }

  return {
    toBottom: toBottom,
    toWindowTop: toWindowTop,
    toWindowBottom: toWindowBottom,
    scrollElement: F2(scrollElement),
    scrollIntoView: scrollIntoView,
    scrollUp: scrollUp,
    scrollDown: scrollDown
  };
}();
