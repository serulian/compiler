$module('casttoanygeneric', function () {
  var $static = this;
  this.$class('48faf528', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.test = function (T) {
    var $f = function (value) {
      $t.cast(value, T, false);
      return;
    };
    return $f;
  };
  $static.TEST = function () {
    $g.casttoanygeneric.test($t.any)($t.fastbox('hello world', $g.____testlib.basictypes.String));
    $g.casttoanygeneric.test($t.any)($t.fastbox(1234, $g.____testlib.basictypes.Integer));
    $g.casttoanygeneric.test($t.any)($g.casttoanygeneric.SomeClass.new());
    $g.casttoanygeneric.test($t.any)($g.____testlib.basictypes.Map($g.____testlib.basictypes.Mappable, $t.any).new());
    $g.casttoanygeneric.test($t.any)($g.____testlib.basictypes.List($t.any).new());
    $g.casttoanygeneric.test($t.any)(null);
    return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
  };
});
