$module('mini', function () {
  var $static = this;
  $static.TEST = function () {
    var lambda;
    lambda = function (someParam) {
      return $t.fastbox(!someParam.$wrapped, $g.________testlib.basictypes.Boolean);
    };
    return lambda($t.fastbox(false, $g.________testlib.basictypes.Boolean));
  };
});
