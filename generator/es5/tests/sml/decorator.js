$module('decorator', function () {
  var $static = this;
  $static.SimpleFunction = function () {
    return $t.fastbox(10, $g.________testlib.basictypes.Integer);
  };
  $static.First = function (decorated, value) {
    return $t.fastbox(decorated.$wrapped + value.$wrapped, $g.________testlib.basictypes.Integer);
  };
  $static.Second = function (decorated, value) {
    return $t.fastbox(decorated.$wrapped - value.$wrapped, $g.________testlib.basictypes.Integer);
  };
  $static.Check = function (decorated, value) {
    return $t.fastbox(value.$wrapped && (decorated.$wrapped == 15), $g.________testlib.basictypes.Boolean);
  };
  $static.TEST = function () {
    return $g.decorator.Check($g.decorator.Second($g.decorator.First($g.decorator.SimpleFunction(), $t.fastbox(10, $g.________testlib.basictypes.Integer)), $t.fastbox(5, $g.________testlib.basictypes.Integer)), $t.fastbox(true, $g.________testlib.basictypes.Boolean));
  };
});
