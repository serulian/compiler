$module('notop', function () {
  var $static = this;
  this.$class('df288f51', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (boolValue) {
      var instance = new $static();
      instance.boolValue = boolValue;
      return instance;
    };
    $static.$bool = function (sc) {
      return sc.boolValue;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var sc;
    sc = $g.notop.SomeClass.new($t.fastbox(false, $g.________testlib.basictypes.Boolean));
    return $t.fastbox(!$g.notop.SomeClass.$bool(sc).$wrapped, $g.________testlib.basictypes.Boolean);
  };
});
