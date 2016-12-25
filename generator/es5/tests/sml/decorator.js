$module('decorator', function () {
  var $static = this;
  $static.SimpleFunction = function () {
    return $t.fastbox(10, $g.____testlib.basictypes.Integer);
  };
  $static.First = function (decorated, value) {
    return $t.fastbox(decorated.$wrapped + value.$wrapped, $g.____testlib.basictypes.Boolean);
  };
  $static.Second = function (decorated, value) {
    return $t.fastbox(decorated.$wrapped - value.$wrapped, $g.____testlib.basictypes.Boolean);
  };
  $static.Check = function (decorated, value) {
    return $t.fastbox(value.$wrapped && (decorated.$wrapped == 15), $g.____testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.decorator.Check($g.decorator.Second($g.decorator.First($g.decorator.SimpleFunction(), $t.fastbox(10, $g.____testlib.basictypes.Integer)), $t.fastbox(5, $g.____testlib.basictypes.Integer)), $t.fastbox(true, $g.____testlib.basictypes.Boolean));
  };
});
