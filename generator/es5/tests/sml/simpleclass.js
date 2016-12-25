$module('simpleclass', function () {
  var $static = this;
  this.$class('f6e326a6', 'SimpleClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $static.Declare = function () {
      return $g.simpleclass.SimpleClass.new();
    };
    $instance.Value = $t.property(function () {
      var $this = this;
      return $t.fastbox(true, $g.____testlib.basictypes.Boolean);
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Declare|1|29dc432d<f6e326a6>": true,
        "Value|3|43834c3f": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var s;
    s = $g.simpleclass.SimpleClass.Declare();
    return s.Value();
  };
});
