$module('requiredfields', function () {
  var $static = this;
  this.$class('4715ab9c', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      instance.SomeField = SomeField;
      instance.AnotherField = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
      return instance;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  $static.TEST = function () {
    var sc;
    sc = $g.requiredfields.SomeClass.new($t.fastbox(2, $g.____testlib.basictypes.Integer));
    return $t.fastbox((sc.SomeField.$wrapped == 2) && sc.AnotherField.$wrapped, $g.____testlib.basictypes.Boolean);
  };
});
