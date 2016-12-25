$module('full', function () {
  var $static = this;
  $static.TEST = function () {
    var lambda;
    lambda = function (firstParam, secondParam) {
      return secondParam;
    };
    return lambda($t.fastbox(123, $g.____testlib.basictypes.Integer), $t.fastbox(true, $g.____testlib.basictypes.Boolean));
  };
});
