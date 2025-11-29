Tea.context(function () {
  this.isUp = "1";

  if (typeof NotifyPopup !== 'undefined') {
    this.success = NotifyPopup;
  } else {
    this.success = function(resp) {
      if (typeof window.parent !== 'undefined' && window.parent.teaweb) {
        window.parent.teaweb.popupFinish(resp);
      }
    }
  }
})