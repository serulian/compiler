$module('memberaccess', function () {
  var $static = this;
  this.$class('a95b5789', 'SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance.someInt = $t.fastbox(2, $g.____testlib.basictypes.Integer);
      instance.someBool = $t.fastbox(true, $g.____testlib.basictypes.Boolean);
      return instance;
    };
    $static.Build = function () {
      return $g.memberaccess.SomeClass.new();
    };
    $instance.InstanceFunc = function () {
      var $this = this;
      return;
    };
    $instance.SomeProp = $t.property(function () {
      var $this = this;
      return $this.someInt;
    });
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Build|1|29dc432d<a95b5789>": true,
        "InstanceFunc|2|29dc432d<void>": true,
        "SomeProp|3|775378f0": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.DoSomething = function (sc, scn) {
    $t.dynamicaccess($g.memberaccess.SomeClass, 'Build', false);
    $t.dynamicaccess($g.memberaccess.SomeClass, 'Build', false);
    $t.dynamicaccess(scn, 'someInt', false);
    $t.dynamicaccess(maimport, 'AnotherFunction', false);
    $g.maimport.AnotherFunction;
    $g.maimport.AnotherFunction;
    sc.InstanceFunc();
    $t.dynamicaccess(sc, 'InstanceFunc', false);
    sc.SomeProp();
    sc.SomeProp();
    scn.SomeProp();
    return;
  };
  $static.TEST = function () {
    var sc;
    sc = $g.memberaccess.SomeClass.new();
    return $t.fastbox(sc.someBool.$wrapped && sc.someBool.$wrapped, $g.____testlib.basictypes.Boolean);
  };
});
