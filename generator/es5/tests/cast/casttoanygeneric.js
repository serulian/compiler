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
    $g.casttoanygeneric.test($t.any)($t.fastbox('hello world', $g.________testlib.basictypes.String));
    $g.casttoanygeneric.test($t.any)($t.fastbox(1234, $g.________testlib.basictypes.Integer));
    $g.casttoanygeneric.test($t.any)($g.casttoanygeneric.SomeClass.new());
    $g.casttoanygeneric.test($t.any)($g.________testlib.basictypes.Mapping($t.any).Empty());
    $g.casttoanygeneric.test($t.any)($g.________testlib.basictypes.Slice($t.any).Empty());
    $g.casttoanygeneric.test($t.any)(null);
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  };
});
