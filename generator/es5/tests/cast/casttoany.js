$module('casttoany', function () {
  var $static = this;
  this.$class('878b6db0', 'SomeClass', false, '', function () {
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

  $static.TEST = function () {
    $t.cast($t.fastbox('hello', $g.________testlib.basictypes.String), $t.any, true);
    $t.cast($t.fastbox(123, $g.________testlib.basictypes.Integer), $t.any, true);
    $t.cast($g.casttoany.SomeClass.new(), $t.any, true);
    $t.cast(null, $t.any, true);
    return $t.fastbox(true, $g.________testlib.basictypes.Boolean);
  };
});
