$module('in', function () {
  var $static = this;
  this.$class('0135b434', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.$contains = function (value) {
      var $this = this;
      return $t.fastbox(!value.$wrapped, $g.____testlib.basictypes.Boolean);
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var sc;
    sc = $g.in.SomeClass.new();
    return sc.$contains($t.fastbox(false, $g.____testlib.basictypes.Boolean));
  };
});
