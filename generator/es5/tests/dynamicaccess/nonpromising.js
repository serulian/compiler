$module('nonpromising', function () {
  var $static = this;
  this.$struct('5d13a931', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (field) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        field: field,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'field', 'field', function () {
      return $g.________testlib.basictypes.Integer;
    }, function () {
      return $g.________testlib.basictypes.Integer;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|fd8bc7c9<5d13a931>": true,
        "equals|4|fd8bc7c9<54ff3ddf>": true,
        "Stringify|2|fd8bc7c9<44e219a9>": true,
        "Mapping|2|fd8bc7c9<ad6de9ce<any>>": true,
        "Clone|2|fd8bc7c9<5d13a931>": true,
        "String|2|fd8bc7c9<44e219a9>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.doSomething = function () {
    var ss;
    var ssa;
    ss = $g.nonpromising.SomeStruct.new($t.fastbox(42, $g.________testlib.basictypes.Integer));
    ssa = ss;
    return $t.cast($t.dynamicaccess(ssa, 'field', false), $g.________testlib.basictypes.Integer, false);
  };
  $static.TEST = function () {
    return $t.fastbox($g.nonpromising.doSomething().$wrapped == 42, $g.________testlib.basictypes.Boolean);
  };
});
