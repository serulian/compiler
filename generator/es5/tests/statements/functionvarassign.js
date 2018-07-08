$module('functionvarassign', function () {
  var $static = this;
  this.$class('e7dcfe11', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (someFunction) {
      var instance = new $static();
      instance.someFunction = someFunction;
      return instance;
    };
    $static.Default = function () {
      return $g.functionvarassign.SomeClass.new(function () {
        return $t.fastbox(41, $g.________testlib.basictypes.Integer);
      });
    };
    $instance.withFunction = function (f) {
      var $this = this;
      $this.someFunction = f;
      return;
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Default|1|6caba86c<e7dcfe11>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var sc;
    sc = $g.functionvarassign.SomeClass.Default();
    sc.withFunction(function () {
      return $t.fastbox(42, $g.________testlib.basictypes.Integer);
    });
    return $t.fastbox(sc.someFunction().$wrapped == 42, $g.________testlib.basictypes.Boolean);
  };
});
