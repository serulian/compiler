$module('slice', function () {
  var $static = this;
  this.$class('5387b00e', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $instance.$slice = function (start, end) {
      var $this = this;
      return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var c;
    c = $g.slice.SomeClass.new();
    c.$slice($t.fastbox(1, $g.____testlib.basictypes.Integer), $t.fastbox(2, $g.____testlib.basictypes.Integer));
    c.$slice(null, $t.fastbox(1, $g.____testlib.basictypes.Integer));
    c.$slice($t.fastbox(1, $g.____testlib.basictypes.Integer), null);
    return c.$slice($t.fastbox(1, $g.____testlib.basictypes.Integer), $t.fastbox(7, $g.____testlib.basictypes.Integer));
  };
});
