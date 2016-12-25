$module('compare', function () {
  var $static = this;
  this.$class('bf8a0308', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $static.$equals = function (first, second) {
      return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    };
    $static.$compare = function (first, second) {
      return $t.fastbox(1, $g.____testlib.basictypes.Integer);
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var first;
    var second;
    first = $g.compare.SomeClass.new();
    second = $g.compare.SomeClass.new();
    $g.compare.SomeClass.$equals(first, second);
    $t.fastbox(!$g.compare.SomeClass.$equals(first, second).$wrapped, $g.____testlib.basictypes.Boolean);
    $t.fastbox($g.compare.SomeClass.$compare(first, second).$wrapped < 0, $g.____testlib.basictypes.Boolean);
    $t.fastbox($g.compare.SomeClass.$compare(first, second).$wrapped > 0, $g.____testlib.basictypes.Boolean);
    $t.fastbox($g.compare.SomeClass.$compare(first, second).$wrapped <= 0, $g.____testlib.basictypes.Boolean);
    $t.fastbox($g.compare.SomeClass.$compare(first, second).$wrapped >= 0, $g.____testlib.basictypes.Boolean);
    return $g.compare.SomeClass.$equals(first, second);
  };
});
